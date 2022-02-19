package update_utils

import (
	"pocok/src/db"
	"pocok/src/services/typless"
	"pocok/src/services/typless/create_training_data"
	"pocok/src/services/wise"
	"pocok/src/utils"
	"pocok/src/utils/models"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type AcceptDependencies struct {
	DbClient       *dynamodb.Client
	TableName      string
	TyplessToken   string
	TyplessDocType string
	WiseService    *wise.WiseService
	WiseQueueUrl   string
	SqsClient      *sqs.Client
}

type AcceptInvoiceInput struct {
	OrgId     string
	InvoiceId string
}

func (d *AcceptDependencies) AcceptInvoice(input AcceptInvoiceInput) error {
	invoice, getInvoiceError := db.GetInvoice(d.DbClient, d.TableName, db.GetInvoiceInput{OrgId: input.OrgId, InvoiceId: input.InvoiceId})
	if getInvoiceError != nil {
		utils.LogError("Invoice query error", getInvoiceError)
		return getInvoiceError
	}

	wiseError := d.queueWiseTransfer(*invoice)
	if wiseError != nil {
		utils.LogError("Wise error", wiseError)
		return wiseError
	}
	updateError := db.UpdateInvoiceStatus(d.DbClient, d.TableName, db.UpdateStatusInput{
		OrgId:     input.OrgId,
		InvoiceId: input.InvoiceId,
		Status:    models.TRANSFER_LOADING,
	})
	if updateError != nil {
		utils.LogError("Update error", updateError)
		return updateError
	}

	typlessError := typless.AddDocumentFeedback(
		&typless.Config{
			Token:   d.TyplessToken,
			DocType: d.TyplessDocType,
		},
		*create_training_data.CreateTrainingData(invoice),
	)
	if typlessError != nil {
		utils.LogError("Error adding document feedback to typless", typlessError)
	}

	return nil
}

func (d *AcceptDependencies) queueWiseTransfer(invoice models.Invoice) error {
	wiseDeps := wise.WiseDependencies{
		WiseService:  d.WiseService,
		SqsClient:    d.SqsClient,
		WiseQueueUrl: d.WiseQueueUrl,
	}

	step0 := &wise.WiseMessageData{
		Invoice: invoice,
	}
	step1, step1Error := wiseDeps.Step1GetProfileId(*step0)
	if step1Error != nil {
		utils.LogError("handler - step1", step1Error)
		return step1Error
	}
	step2, step2Error := wiseDeps.Step2UpsertRecipientAccount(*step1)
	if step2Error != nil {
		utils.LogError("handler - step2", step2Error)
		return step2Error
	}
	step3, step3Error := wiseDeps.Step3CreateQuote(*step2)
	if step3Error != nil {
		utils.LogError("handler - step3", step3Error)
		return step3Error
	}
	step3MessageError := wiseDeps.SendMessage(step3)
	if step3MessageError != nil {
		utils.LogError(step3MessageError.Error(), step3MessageError)
		return step3MessageError
	}
	// The flow continues in the wise transfer consumer
	return nil
}
