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
		return utils.MailApiResponse(http.StatusUnauthorized, ""), parseTokenError
	}

	invoices, getPendingInvoicesError := db.GetPendingInvoices(d.dbClient, d.tableName, claims.OrgId)
	if getPendingInvoicesError != nil {
		utils.LogError("Error while getting pending invoices from db", getPendingInvoicesError)
		return nil, getPendingInvoicesError
	}

	invoicesWithLinks, presignError := getInvoicesWithLinks(d, invoices)
	if presignError != nil {
		utils.LogError(presignError.Error(), presignError)
		return nil, presignError
	}
	response := models.InvoiceResponse{
		Items: invoicesWithLinks,
		Total: len(invoices),
	}

	invoiceBytes, marshalError := json.Marshal(response)
	if marshalError != nil {
		utils.LogError("Error while parsing invoices from db", marshalError)
		return nil, marshalError
	}

	invoiceStr := string(invoiceBytes)
	return utils.MailApiResponse(http.StatusOK, invoiceStr), nil
}

func getInvoicesWithLinks(d *dependencies, invoices []models.Invoice) ([]models.InvoiceWithLink, error) {
	invoicesWithLinks := make([]models.InvoiceWithLink, len(invoices))
	psClient := s3.NewPresignClient(d.s3Client)
	for i, invoice := range invoices {
		input := &s3.GetObjectInput{
			Bucket: &d.bucketName,
			Key:    &invoice.Filename,
		}
		resp, presignError := psClient.PresignGetObject(context.TODO(), input, s3.WithPresignExpires(time.Hour))
		if presignError != nil {
			utils.LogError(presignError.Error(), presignError)
			return nil, presignError
		}

		invoicesWithLinks[i] = models.InvoiceWithLink{
			Invoice: invoice,
			Link:    resp.URL,
		}
	}
	return invoicesWithLinks, nil
}
