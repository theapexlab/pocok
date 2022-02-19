package main

import (
	"errors"
	"net/http"
	"os"
	"pocok/src/api/invoices/update_utils"
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
		return utils.MailApiResponse(http.StatusUnauthorized, utils.ApiErrorBody(parseTokenError.Error())), nil
	}

	data, parseFormDataError := request_parser.ParseUrlEncodedFormData(r)
	if parseFormDataError != nil {
		utils.LogError("Form body parse failed", parseFormDataError)
		return utils.MailApiResponse(http.StatusBadRequest, utils.ApiErrorBody(parseFormDataError.Error())), nil
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
		acceptError := deps.AcceptInvoice(update_utils.AcceptInvoiceInput{
			OrgId:     claims.OrgId,
			InvoiceId: invoiceId,
		})
		if acceptError != nil {
			utils.LogError("error while accepting invoice", acceptError)
			acceptErrors += acceptError.Error()
		}
	}

	if len(acceptErrors) != 0 {
		acceptErrorSummary := errors.New(acceptErrors)
		utils.LogError("error while accepting invoices", acceptErrorSummary)
		return utils.MailApiResponse(http.StatusInternalServerError, acceptErrorSummary.Error()), nil
	}
	return utils.MailApiResponse(http.StatusOK, ""), nil
}
