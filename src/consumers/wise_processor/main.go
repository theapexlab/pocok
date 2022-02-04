package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"pocok/src/services/wise"
	apiModels "pocok/src/services/wise/api/models"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
)

type dependencies struct {
	wiseQueueUrl string
	sqsClient    *sqs.Client
	wiseService  *wise.WiseService
}

func main() {
	d := &dependencies{
		wiseQueueUrl: os.Getenv("queueUrl"),
		sqsClient:    aws_clients.GetSQSClient(),
		wiseService:  wise.CreateWiseService(os.Getenv("wiseApiToken")),
	}

	lambda.Start(d.handler)
}

func (d *dependencies) handler(event events.SQSEvent) error {
	for _, record := range event.Records {
		var messageData wise.WiseMessageData
		if unmarshalError := json.Unmarshal([]byte(record.Body), &messageData); unmarshalError != nil {
			utils.LogError("handler - Unmarshal", unmarshalError)
			return unmarshalError
		}

		switch eventType := messageData.RequestType; eventType {
		case wise.WiseStep1:
			newMessage, step1Error := d.step1GetProfileId(messageData)
			if step1Error != nil {
				utils.LogError("handler - step1", step1Error)
				return step1Error
			}
			step1MessageError := d.sendMessage(newMessage)
			if step1MessageError != nil {
				utils.LogError(step1MessageError.Error(), step1MessageError)
				return step1MessageError
			}
		case wise.WiseStep2:
			newMessage, step2Error := d.step2UpsertRecipientAccount(messageData)
			if step2Error != nil {
				utils.LogError("handler - step2", step2Error)
				return step2Error
			}
			step2MessageError := d.sendMessage(newMessage)
			if step2MessageError != nil {
				utils.LogError(step2MessageError.Error(), step2MessageError)
				return step2MessageError
			}
		case wise.WiseStep3:
			newMessage, step3Error := d.step3CreateQuote(messageData)
			if step3Error != nil {
				utils.LogError("handler - step3", step3Error)
				return step3Error
			}
			step3MessageError := d.sendMessage(newMessage)
			if step3MessageError != nil {
				utils.LogError(step3MessageError.Error(), step3MessageError)
				return step3MessageError
			}
		case wise.WiseStep4:
			step4Error := d.step4CreateTransfer(messageData)
			if step4Error != nil {
				utils.LogError("handler - step4", step4Error)
				return step4Error
			}
		default:
			return errors.New("unknown event type")
		}
	}
	return nil
}

func (d *dependencies) step1GetProfileId(step1Data wise.WiseMessageData) (*wise.WiseMessageData, error) {
	profile, getBusinessProfileError := d.wiseService.GetBusinessProfile()
	if getBusinessProfileError != nil {
		utils.LogError("step1GetProfileId - GetBusinessProfile", getBusinessProfileError)
		return nil, getBusinessProfileError
	}
	step2Data := step1Data
	step2Data.RequestType = wise.WiseStep2
	step2Data.ProfileId = profile.ID
	return &step2Data, nil
}

func (d *dependencies) step2UpsertRecipientAccount(step2Data wise.WiseMessageData) (*wise.WiseMessageData, error) {
	recipient, upsertRecipientAccountError := d.wiseService.UpsertRecipient(&step2Data.Invoice)
	if upsertRecipientAccountError != nil {
		utils.LogError("step2UpsertRecipientAccount - UpsertRecipient", upsertRecipientAccountError)
		return nil, upsertRecipientAccountError
	}
	step3Data := step2Data
	step3Data.RequestType = wise.WiseStep3
	step3Data.RecipientAccountId = recipient.ID
	return &step3Data, nil
}

func (d *dependencies) step3CreateQuote(step3Data wise.WiseMessageData) (*wise.WiseMessageData, error) {
	grossPrice, atoiError := strconv.Atoi(step3Data.Invoice.GrossPrice)
	if atoiError != nil {
		utils.LogError("step3CreateQuote - Atoi", atoiError)
		return nil, atoiError
	}
	quoteInput := apiModels.Quote{
		Profile:        step3Data.ProfileId,
		TargetAccount:  step3Data.RecipientAccountId,
		PayOut:         "BALANCE",
		SourceCurrency: "HUF",
		TargetCurrency: step3Data.Invoice.Currency,
		TargetAmount:   grossPrice,
	}
	quote, createQuoteError := d.wiseService.WiseApi.CreateQuote(quoteInput)
	if createQuoteError != nil {
		utils.LogError("step3CreateQuote - CreateQuote", createQuoteError)
		return nil, createQuoteError
	}
	step4Data := step3Data
	step4Data.RequestType = wise.WiseStep4
	step4Data.QuoteId = quote.ID
	step4Data.TransactionId = uuid.NewString()

	return &step4Data, nil
}

func (d *dependencies) step4CreateTransfer(step4Data wise.WiseMessageData) error {
	transferInput := apiModels.Transfer{
		TargetAccount: step4Data.RecipientAccountId,
		QuoteUUID:     step4Data.QuoteId,
		Details: struct {
			Reference string `json:"reference"`
		}{Reference: step4Data.Invoice.InvoiceNumber},
		CustomerTransactionID: step4Data.TransactionId,
	}
	_, createTransferError := d.wiseService.WiseApi.CreateTransfer(transferInput)
	if createTransferError != nil {
		utils.LogError("step4CreateTransfer - CreateTransfer", createTransferError)
		return createTransferError
	}

	return nil
}

func (d *dependencies) sendMessage(wiseMessage *wise.WiseMessageData) error {
	messageByteArray, marshalError := json.Marshal(wiseMessage)
	if marshalError != nil {
		utils.LogError("sendMessage - Marshal", marshalError)
		return marshalError
	}
	fmt.Println(d.wiseQueueUrl)
	fmt.Println(wiseMessage)
	messageString := string(messageByteArray)
	sqsResponse, sqsError := d.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: &messageString,
		QueueUrl:    &d.wiseQueueUrl,
	})
	fmt.Println(sqsResponse)
	fmt.Println(sqsError)
	return sqsError
}
