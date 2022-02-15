package main

import (
	"errors"
	"net/http"
	"os"
	"pocok/src/api/invoices/update_utils"
	"pocok/src/db"
	"pocok/src/services/wise"
	"pocok/src/utils"
	"pocok/src/utils/auth"
	"pocok/src/utils/aws_clients"
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
	wiseService    *wise.WiseService
	wiseQueueUrl   string
	sqsClient      *sqs.Client
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

	deps := update_utils.AcceptDependencies{
		DbClient:       d.dbClient,
		TableName:      d.tableName,
		TyplessToken:   d.typlessToken,
		TyplessDocType: d.typlessDocType,
		WiseService:    d.wiseService,
		WiseQueueUrl:   d.wiseQueueUrl,
		SqsClient:      d.sqsClient,
	}
	acceptErrors := ""
	for invoiceId := range data {
		acceptError := deps.AcceptInvoice(*claims, db.StatusUpdate{
			InvoiceId: invoiceId,
		})
		if acceptError != nil {
			acceptErrors += acceptError.Error()
		}
	}

	if len(acceptErrors) != 0 {
		return utils.MailApiResponse(http.StatusInternalServerError, "{}"), errors.New(acceptErrors)
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
