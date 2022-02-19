package main

import (
	"errors"
	"net/http"
	"os"
	"pocok/src/db"
	"pocok/src/utils"
	"pocok/src/utils/auth"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/request_parser"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type dependencies struct {
	dbClient   *dynamodb.Client
	tableName  string
	s3Client   *s3.Client
	bucketName string
}

type formData struct {
	InvoiceId string
}

func main() {
	d := &dependencies{
		dbClient:   aws_clients.GetDbClient(),
		tableName:  os.Getenv("tableName"),
		s3Client:   aws_clients.GetS3Client(),
		bucketName: os.Getenv("bucketName"),
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

	deleteError := db.DeleteInvoice(d.dbClient, d.tableName, *d.s3Client, d.bucketName, db.DeleteInvoiceInput{OrgId: claims.OrgId, InvoiceId: data.InvoiceId})
	if deleteError != nil {
		utils.LogError("Error while removing invoice", deleteError)
		return utils.MailApiResponse(http.StatusInternalServerError, utils.ApiErrorBody(deleteError.Error())), nil
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
