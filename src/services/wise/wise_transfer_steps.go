package wise

import (
	"context"
	"encoding/json"
	apiModels "pocok/src/services/wise/api/models"
	"pocok/src/utils"
	"pocok/src/utils/currency"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
)

type WiseDependencies struct {
	WiseService  *WiseService
	SqsClient    *sqs.Client
	WiseQueueUrl string
}

func (d *WiseDependencies) Step1GetProfileId(step1Data WiseMessageData) (*WiseMessageData, error) {
	profile, getBusinessProfileError := d.WiseService.GetBusinessProfile()
	if getBusinessProfileError != nil {
		utils.LogError("step1GetProfileId - GetBusinessProfile", getBusinessProfileError)
		return nil, getBusinessProfileError
	}
	step2Data := step1Data
	step2Data.RequestType = WiseStep2
	step2Data.ProfileId = profile.ID
	return &step2Data, nil
}

func (d *WiseDependencies) Step2UpsertRecipientAccount(step2Data WiseMessageData) (*WiseMessageData, error) {
	recipient, upsertRecipientAccountError := d.WiseService.UpsertRecipient(step2Data.ProfileId, &step2Data.Invoice)
	if upsertRecipientAccountError != nil {
		utils.LogError("step2UpsertRecipientAccount - UpsertRecipient", upsertRecipientAccountError)
		return nil, upsertRecipientAccountError
	}
	step3Data := step2Data
	step3Data.RequestType = WiseStep3
	step3Data.RecipientAccountId = recipient.ID
	return &step3Data, nil
}

func (d *WiseDependencies) Step3CreateQuote(step3Data WiseMessageData) (*WiseMessageData, error) {
	grossPriceString := currency.GetValueFromPrice(step3Data.Invoice.GrossPrice)
	grossPrice, parseFloatError := strconv.ParseFloat(grossPriceString, 32)
	if parseFloatError != nil {
		utils.LogError("step3CreateQuote - Parse float", parseFloatError)
		return nil, parseFloatError
	}
	quoteInput := apiModels.Quote{
		Profile:        step3Data.ProfileId,
		TargetAccount:  step3Data.RecipientAccountId,
		PayOut:         "BALANCE",
		SourceCurrency: "HUF",
		TargetCurrency: step3Data.Invoice.Currency,
		TargetAmount:   grossPrice,
	}
	quote, createQuoteError := d.WiseService.WiseApi.CreateQuote(quoteInput)
	if createQuoteError != nil {
		utils.LogError("step3CreateQuote - CreateQuote", createQuoteError)
		return nil, createQuoteError
	}
	step4Data := step3Data
	step4Data.RequestType = WiseStep4
	step4Data.QuoteId = quote.ID
	step4Data.TransactionId = uuid.NewString()

	return &step4Data, nil
}

func (d *WiseDependencies) Step4CreateTransfer(step4Data WiseMessageData) error {
	transferInput := apiModels.Transfer{
		TargetAccount: step4Data.RecipientAccountId,
		QuoteUUID:     step4Data.QuoteId,
		Details: struct {
			Reference string `json:"reference"`
		}{Reference: step4Data.Invoice.InvoiceNumber},
		CustomerTransactionID: step4Data.TransactionId,
	}
	_, createTransferError := d.WiseService.WiseApi.CreateTransfer(transferInput)
	if createTransferError != nil {
		utils.LogError("step4CreateTransfer - CreateTransfer", createTransferError)
		return createTransferError
	}

	return nil
}

func (d *WiseDependencies) SendMessage(wiseMessage *WiseMessageData) error {
	messageByteArray, marshalError := json.Marshal(wiseMessage)
	if marshalError != nil {
		utils.LogError("sendMessage - Marshal", marshalError)
		return marshalError
	}
	messageString := string(messageByteArray)
	_, sqsError := d.SqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: &messageString,
		QueueUrl:    &d.WiseQueueUrl,
	})
	return sqsError
}
