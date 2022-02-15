package update_utils

import (
	"pocok/src/db"
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

func (d *AcceptDependencies) AcceptInvoice(claims models.JWTClaims, updateInput db.StatusUpdate) error {
	invoice, getInvoiceError := db.GetInvoice(d.DbClient, d.TableName, claims.OrgId, updateInput.InvoiceId)
	if getInvoiceError != nil {
		utils.LogError("Invoice query error", getInvoiceError)
		return getInvoiceError
	}
	wiseDeps := WiseDependencies{
		WiseService:  d.WiseService,
		SqsClient:    d.SqsClient,
		WiseQueueUrl: d.WiseQueueUrl,
	}
	wiseError := wiseDeps.WiseSteps(*invoice)
	if wiseError != nil {
		utils.LogError("Wise error", wiseError)
		return wiseError
	}
	updateError := db.UpdateInvoiceStatus(d.DbClient, d.TableName, claims.OrgId, db.StatusUpdate{
		InvoiceId: updateInput.InvoiceId,
		Status:    models.TRANSFER_LOADING,
	})
	if updateError != nil {
		utils.LogError("Update error", updateError)
		return updateError
	}
	typlessError := UpdateTypless(d.TyplessToken, d.TyplessDocType, *invoice)
	if typlessError != nil {
		utils.LogError("Error while submitting typless feedback", typlessError)
	}
	return nil
}
