package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"pocok/src/db"
	"pocok/src/utils"
	"pocok/src/utils/auth"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type dependencies struct {
	dbClient   *dynamodb.Client
	tableName  string
	bucketName string
	s3Client   *s3.Client
}

func main() {
	d := &dependencies{
		tableName:  os.Getenv("tableName"),
		dbClient:   aws_clients.GetDbClient(),
		bucketName: os.Getenv("bucketName"),
		s3Client:   aws_clients.GetS3Client(),
	}

	lambda.Start(d.handler)
}

func (d *dependencies) handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	token := request.QueryStringParameters["token"]
	claims, parseTokenError := auth.ParseToken(token)
	if parseTokenError != nil {
		return utils.MailApiResponse(http.StatusUnauthorized, utils.ApiErrorBody(parseTokenError.Error())), nil
	}

	invoices, getPendingInvoicesError := db.GetPendingInvoices(d.dbClient, d.tableName, claims.OrgId)
	if getPendingInvoicesError != nil {
		utils.LogError("Error while getting pending invoices from db", getPendingInvoicesError)
		return utils.MailApiResponse(http.StatusInternalServerError, utils.ApiErrorBody(getPendingInvoicesError.Error())), nil

	}

	response, responseError := d.getInvoiceResponse(invoices)
	if responseError != nil {
		utils.LogError(responseError.Error(), responseError)
		return utils.MailApiResponse(http.StatusInternalServerError, utils.ApiErrorBody(responseError.Error())), nil

	}

	invoiceBytes, marshalError := json.Marshal(response)
	if marshalError != nil {
		utils.LogError("Error while parsing invoices from db", marshalError)
		return utils.MailApiResponse(http.StatusInternalServerError, marshalError.Error()), nil

	}

	invoiceStr := string(invoiceBytes)
	return utils.MailApiResponse(http.StatusOK, invoiceStr), nil
}

func (d *dependencies) getInvoiceResponse(invoices []models.Invoice) (*models.InvoiceResponse, error) {
	psClient := s3.NewPresignClient(d.s3Client)

	items := make([]models.InvoiceResponseItem, len(invoices))
	for i, invoice := range invoices {
		extendedInvoice := utils.ExtendInvoice(invoice)

		input := &s3.GetObjectInput{
			Bucket: &d.bucketName,
			Key:    &invoice.Filename,
		}
		resp, presignError := psClient.PresignGetObject(context.TODO(), input, s3.WithPresignExpires(time.Hour))
		if presignError != nil {
			utils.LogError(presignError.Error(), presignError)
			return nil, presignError
		}

		items[i] = models.InvoiceResponseItem{
			Invoice: *extendedInvoice,
			Link:    resp.URL,
		}
	}
	response := models.InvoiceResponse{
		Items: items,
		Total: len(items),
	}

	return &response, nil
}
