package main

import (
	"net/http"
	"os"
	"pocok/src/api/invoices/update_utils"
	"pocok/src/db"
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
}

func main() {
	d := &dependencies{
		tableName: os.Getenv("tableName"),
		dbClient:  aws_clients.GetDbClient(),
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

	var invoiceIds []string
	for invoiceId := range data {
		invoiceIds = append(invoiceIds, invoiceId)
	}

	updateError := db.UpdateInvoiceStatuses(d.dbClient, d.tableName, claims.OrgId, invoiceIds, models.ACCEPTED)
	if updateError != nil {
		utils.LogError("Error updating dynamo db", updateError)
		return utils.MailApiResponse(http.StatusInternalServerError, ""), nil
	}

	for invoiceId := range data {
		invoice, getInvoiceError := db.GetInvoice(d.dbClient, d.tableName, claims.OrgId, invoiceId)
		if getInvoiceError != nil {
			utils.LogError("Error while getting invoice", getInvoiceError)
			continue
		}

		feedbackError := update_utils.UpdateTypless(d.typlessToken, d.typlessDocType, *invoice)
		if feedbackError != nil {
			utils.LogError("Error while submitting typless feedback", feedbackError)
		}

		wiseError := update_utils.SendWiseMessage(*d.sqsClient, d.wiseQueueUrl, *invoice)
		if wiseError != nil {
			utils.LogError("Error while creating wise request", wiseError)
		}
	}

	return utils.MailApiResponse(http.StatusOK, "{}"), nil
}
