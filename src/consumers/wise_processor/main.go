package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"pocok/src/services/wise"
	apiModels "pocok/src/services/wise/api/models"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
)

const (
	step1 = "step1:get_profile_id"
	step2 = "step2:upsert_recipient_account"
	step3 = "step3:create_quote"
	step4 = "step4:create_transfer"
)

type wiseMessageData struct {
	RequestType        string         `json:"requestType"`
	ProfileId          int            `json:"profileId"`
	RecipientAccountId int            `json:"recipientAccountId"`
	QuoteId            string         `json:"quoteId"`
	TransactionId      string         `json:"transactionId"`
	Reference          string         `json:"reference"`
	Invoice            models.Invoice `json:"invoice"`
}

type dependencies struct {
	wiseQueueUrl string
	sqsClient    *sqs.Client
	wiseService  *wise.WiseService
}

func main() {
	d := &dependencies{
		wiseQueueUrl: os.Getenv("wiseQueueUrl"),
		sqsClient:    aws_clients.GetSQSClient(),
		wiseService:  wise.CreateWiseService(os.Getenv("")),
	}

	lambda.Start(d.handler)
}

func (d *dependencies) handler(event events.SQSEvent) error {
	for _, record := range event.Records {
		var messageData wiseMessageData
		if unmarshalError := json.Unmarshal([]byte(record.Body), &messageData); unmarshalError != nil {
			utils.LogError("handler - Unmarshal", unmarshalError)
			return unmarshalError
		}
		switch eventType := messageData.RequestType; eventType {
		case step1:
			newMessage, step1Error := d.step1GetProfileId(messageData)
			if step1Error != nil {
				utils.LogError("handler - step1", step1Error)
				return step1Error
			}
			d.sendMessage(newMessage)
		case step2:
			newMessage, step2Error := d.step2UpsertRecipientAccount(messageData)
			if step2Error != nil {
				utils.LogError("handler - step2", step2Error)
				return step2Error
			}
			d.sendMessage(newMessage)
		case step3:
			newMessage, step3Error := d.step3CreateQuote(messageData)
			if step3Error != nil {
				utils.LogError("handler - step3", step3Error)
				return step3Error
			}
			d.sendMessage(newMessage)
		case step4:
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

func (d *dependencies) step1GetProfileId(step1Data wiseMessageData) (*wiseMessageData, error) {
	profile, getBusinessProfileError := d.wiseService.GetBusinessProfile()
	if getBusinessProfileError != nil {
		utils.LogError("step1GetProfileId - GetBusinessProfile", getBusinessProfileError)
		return nil, getBusinessProfileError
	}
	step2Data := step1Data
	step2Data.RequestType = step2
	step2Data.ProfileId = profile.ID
	return &step2Data, nil
}

func (d *dependencies) step2UpsertRecipientAccount(step2Data wiseMessageData) (*wiseMessageData, error) {
	recipient, upsertRecipientAccountError := d.wiseService.UpsertRecipient(&step2Data.Invoice)
	if upsertRecipientAccountError != nil {
		utils.LogError("step2UpsertRecipientAccount - UpsertRecipient", upsertRecipientAccountError)
		return nil, upsertRecipientAccountError
	}
	step3Data := step2Data
	step3Data.RequestType = step3
	step3Data.RecipientAccountId = recipient.ID
	return &step3Data, nil
}

func (d *dependencies) step3CreateQuote(step3Data wiseMessageData) (*wiseMessageData, error) {
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
	quote, createQuoteError := d.wiseService.CreateQuote(quoteInput)
	if createQuoteError != nil {
		utils.LogError("step3CreateQuote - CreateQuote", createQuoteError)
		return nil, createQuoteError
	}
	step4Data := step3Data
	step4Data.RequestType = step4
	step4Data.QuoteId = quote.ID
	step4Data.TransactionId = uuid.NewString()

	return &step4Data, nil
}

func (d *dependencies) step4CreateTransfer(step4Data wiseMessageData) error {
	transferInput := apiModels.Transfer{
		TargetAccount: step4Data.RecipientAccountId,
		QuoteUUID:     step4Data.QuoteId,
		Details: struct {
			Reference string `json:"reference"`
		}{Reference: step4Data.Invoice.InvoiceNumber},
		CustomerTransactionID: step4Data.TransactionId,
	}
	_, createTransferError := d.wiseService.CreateTransfer(transferInput)
	if createTransferError != nil {
		utils.LogError("step4CreateTransfer - CreateTransfer", createTransferError)
		return createTransferError
	}

	return nil
}

func (d *dependencies) sendMessage(wiseMessage *wiseMessageData) error {
	messageByteArray, marshalError := json.Marshal(wiseMessage)
	if marshalError != nil {
		utils.LogError("sendMessage - Marshal", marshalError)
		return marshalError
	}
	messageString := string(messageByteArray)
	_, sqsError := d.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: &messageString,
		QueueUrl:    &d.wiseQueueUrl,
	})
	return sqsError
}
