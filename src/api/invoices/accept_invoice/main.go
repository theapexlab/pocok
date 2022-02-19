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
	wiseQueueUrl   string
	sqsClient      *sqs.Client
	wiseService    *wise.WiseService
}

type formData struct {
	InvoiceId string
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

	data, parseError := getRequestData(r)
	if parseError != nil {
		utils.LogError("Form body parse failed", parseError)
		return utils.MailApiResponse(http.StatusUnprocessableEntity, utils.ApiErrorBody(parseError.Error())), nil
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
	acceptError := deps.AcceptInvoice(update_utils.AcceptInvoiceInput{
		OrgId:     claims.OrgId,
		InvoiceId: data.InvoiceId,
	})

	if acceptError != nil {
		utils.LogError("Invoice accept failed", acceptError)
		return utils.MailApiResponse(http.StatusInternalServerError, utils.ApiErrorBody(acceptError.Error())), nil
	}

	return utils.MailApiResponse(http.StatusOK, ""), nil
}

func getRequestData(r events.APIGatewayProxyRequest) (*formData, error) {
	mapData, parseFormDataError := request_parser.ParseUrlEncodedFormData(r)
	if parseFormDataError != nil {
		utils.LogError("Form body parse failed", parseFormDataError)
		return nil, parseFormDataError
	}

	var data formData
	mapError := utils.MapToStruct(mapData, &data)
	if mapError != nil {
		return nil, errors.New("invalid input")
	}
	return &data, nil
}
