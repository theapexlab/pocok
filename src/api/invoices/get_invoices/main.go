package main

import (
	"encoding/json"
	"net/http"
	"os"
	"pocok/src/db"
	"pocok/src/utils"
	"pocok/src/utils/auth"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"

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

func (d *dependencies) handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	token := request.QueryStringParameters["token"]
	claims, err := auth.ParseToken(token)
	if err != nil {
		return utils.MailApiResponse(http.StatusUnauthorized, ""), err
	}

	invoices, err := db.GetPendingInvoices(d.dbClient, d.tableName, claims.OrgId)
	if err != nil {
		utils.LogError("Error while getting pending invoices from db", err)
		return nil, err
	}

	// For testing use:
	// invoices := mocks.Invoices

	indexedInvoices := utils.MapInvoiceToInvoiceServiceIndexes(invoices)

	response := models.InvoiceResponse{
		Items: indexedInvoices,
		Total: len(invoices),
	}

	invoiceBytes, err := json.Marshal(response)
	if err != nil {
		utils.LogError("Error while parsing invoices from db", err)
		return nil, err
	}

	invoiceStr := string(invoiceBytes)
	return utils.MailApiResponse(http.StatusOK, invoiceStr), nil
}
