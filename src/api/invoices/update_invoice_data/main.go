package main

import (
	"errors"
	"net/http"
	"os"
	"pocok/src/db"
	"pocok/src/utils"
	"pocok/src/utils/auth"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"
	"pocok/src/utils/request_parser"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type dependencies struct {
	dbClient  *dynamodb.Client
	tableName string
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
		return utils.MailApiResponse(http.StatusUnauthorized, utils.ApiErrorBody(parseTokenError.Error())), nil
	}

	invoice, parseError := getRequestData(r)
	if parseError != nil {
		utils.LogError("Form body parse failed", parseError)
		return utils.MailApiResponse(http.StatusUnprocessableEntity, utils.ApiErrorBody(parseError.Error())), nil
	}

	updateErr := db.UpdateInvoiceData(d.dbClient, d.tableName, db.UpdateDataInput{OrgId: claims.OrgId, Invoice: *invoice})
	if updateErr != nil {
		utils.LogError("Error updating dynamo db", updateErr)
		return utils.MailApiResponse(http.StatusInternalServerError, utils.ApiErrorBody(parseError.Error())), nil
	}

	return utils.MailApiResponse(http.StatusOK, ""), nil
}

func getRequestData(r events.APIGatewayProxyRequest) (*models.Invoice, error) {
	mapData, parseFormDataError := request_parser.ParseUrlEncodedFormData(r)
	if parseFormDataError != nil {
		utils.LogError("Form body parse failed", parseFormDataError)
		return nil, parseFormDataError
	}

	var data models.Invoice
	mapError := utils.MapToStruct(mapData, &data)
	if mapError != nil {
		return nil, errors.New("invalid input")
	}
	return &data, nil
}
