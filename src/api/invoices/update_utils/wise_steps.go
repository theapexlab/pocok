package update_utils

import (
	"context"
	"encoding/json"
	"pocok/src/services/wise"
	apiModels "pocok/src/services/wise/api/models"
	"pocok/src/utils"
	"pocok/src/utils/currency"
	"pocok/src/utils/models"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
)

type WiseDependencies struct {
	WiseService  *wise.WiseService
	SqsClient    *sqs.Client
	WiseQueueUrl string
}

func (d *WiseDependencies) WiseSteps(invoice models.Invoice) error {
	step0 := &wise.WiseMessageData{
		Invoice: invoice,
	}
	step1, step1Error := d.step1GetProfileId(*step0)
	if step1Error != nil {
		utils.LogError("handler - step1", step1Error)
		return step1Error
	}
	step2, step2Error := d.step2UpsertRecipientAccount(*step1)
	if step2Error != nil {
		utils.LogError("handler - step2", step2Error)
		return step2Error
	}
	step3, step3Error := d.step3CreateQuote(*step2)
	if step3Error != nil {
		utils.LogError("handler - step3", step3Error)
		return step3Error
	}
	// The flow continues in the wise transfer consumer
	step3MessageError := d.sendMessage(step3)
	if step3MessageError != nil {
		utils.LogError(step3MessageError.Error(), step3MessageError)
		return step3MessageError
	}
	return nil
}

func (d *WiseDependencies) step1GetProfileId(step1Data wise.WiseMessageData) (*wise.WiseMessageData, error) {
	profile, getBusinessProfileError := d.WiseService.GetBusinessProfile()
	if getBusinessProfileError != nil {
		utils.LogError("step1GetProfileId - GetBusinessProfile", getBusinessProfileError)
		return nil, getBusinessProfileError
	}
	step2Data := step1Data
	step2Data.RequestType = wise.WiseStep2
	step2Data.ProfileId = profile.ID
	return &step2Data, nil
}

func (d *WiseDependencies) step2UpsertRecipientAccount(step2Data wise.WiseMessageData) (*wise.WiseMessageData, error) {
	recipient, upsertRecipientAccountError := d.WiseService.UpsertRecipient(step2Data.ProfileId, &step2Data.Invoice)
	if upsertRecipientAccountError != nil {
		utils.LogError("step2UpsertRecipientAccount - UpsertRecipient", upsertRecipientAccountError)
		return nil, upsertRecipientAccountError
	}
	step3Data := step2Data
	step3Data.RequestType = wise.WiseStep3
	step3Data.RecipientAccountId = recipient.ID
	return &step3Data, nil
}

func (d *WiseDependencies) step3CreateQuote(step3Data wise.WiseMessageData) (*wise.WiseMessageData, error) {
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
	step4Data.RequestType = wise.WiseStep4
	step4Data.QuoteId = quote.ID
	step4Data.TransactionId = uuid.NewString()

	return &step4Data, nil
}

/*
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
*/

func (d *WiseDependencies) sendMessage(wiseMessage *wise.WiseMessageData) error {
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
