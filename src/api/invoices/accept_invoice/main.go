package main

import (
	"net/http"
	"os"
	"pocok/src/api/invoices/update_utils"
	"pocok/src/db"
	"pocok/src/services/wise"
	"pocok/src/utils"
	"pocok/src/utils/auth"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"
	"pocok/src/utils/request_parser"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type dependencies struct {
	dbClient       *dynamodb.Client
	tableName      string
	typlessToken   string
	typlessDocType string
	wiseQueueUrl   string
	sqsClient      *sqs.Client
	wiseService    *wise.WiseService
}

func main() {
	d := &dependencies{
		dbClient:       aws_clients.GetDbClient(),
		tableName:      os.Getenv("tableName"),
		typlessToken:   os.Getenv("typlessToken"),
		typlessDocType: os.Getenv("typlessDocType"),
		wiseQueueUrl:   os.Getenv("wiseQueueUrl"),
		sqsClient:      aws_clients.GetSQSClient(),
		wiseService:    wise.CreateWiseService(os.Getenv("wiseApiToken")),
	}
	lambda.Start(d.handler)
}

func (d *dependencies) handler(r events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	token := r.QueryStringParameters["token"]
	claims, parseTokenError := auth.ParseToken(token)

	if parseTokenError != nil {
		utils.LogError("Token validation failed", parseTokenError)
		return utils.MailApiResponse(http.StatusUnauthorized, ""), parseTokenError
	}

	data, parseFormDataError := request_parser.ParseUrlEncodedFormData(r)
	if parseFormDataError != nil {
		utils.LogError("Form body parse failed", parseFormDataError)
		return utils.MailApiResponse(http.StatusBadRequest, ""), parseFormDataError
	}

	update, validationError := db.GetValidStatusUpdate(data)
	if validationError != nil {
		utils.LogError("Invalid while validating update", validationError)
		return utils.MailApiResponse(http.StatusUnprocessableEntity, validationError.Error()), validationError
	}

	acceptError := d.AcceptInvoice(*claims, *update)
	if acceptError != nil {
		utils.LogError("Invoice accept failed", validationError)
		return utils.MailApiResponse(http.StatusInternalServerError, acceptError.Error()), acceptError
	}

	return utils.MailApiResponse(http.StatusOK, ""), nil
}

func (d *dependencies) AcceptInvoice(claims models.JWTClaims, updateInput db.StatusUpdate) error {
	invoice, getInvoiceError := db.GetInvoice(d.dbClient, d.tableName, claims.OrgId, updateInput.InvoiceId)
	if getInvoiceError != nil {
		utils.LogError("Invoice query error", getInvoiceError)
		return getInvoiceError
	}
	wiseDeps := update_utils.WiseDependencies{
		WiseService:  d.wiseService,
		SqsClient:    d.sqsClient,
		WiseQueueUrl: d.wiseQueueUrl,
	}
	wiseError := wiseDeps.WiseSteps(*invoice)
	if wiseError != nil {
		utils.LogError("Wise error", wiseError)
		return wiseError
	}
	updateError := db.UpdateInvoiceStatus(d.dbClient, d.tableName, claims.OrgId, db.StatusUpdate{
		InvoiceId: updateInput.InvoiceId,
		Status:    models.TRANSFER_LOADING,
	})
	if updateError != nil {
		utils.LogError("Update error", updateError)
		return updateError
	}
	typlessError := update_utils.UpdateTypless(d.typlessToken, d.typlessDocType, *invoice)
	if typlessError != nil {
		utils.LogError("Error while submitting typless feedback", typlessError)
	}
	return nil
}
