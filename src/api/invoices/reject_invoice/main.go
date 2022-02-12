package main

import (
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

	deleteError := db.DeleteInvoice(d.dbClient, d.tableName, *d.s3Client, d.bucketName, claims.OrgId, update.InvoiceId, update.Filename)
	if deleteError != nil {
		utils.LogError("Error while removing invoice", deleteError)
		return utils.MailApiResponse(http.StatusInternalServerError, ""), deleteError
	}

	return utils.MailApiResponse(http.StatusOK, ""), nil
}
