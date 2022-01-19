package main

import (
	"encoding/json"
	"net/http"
	"os"
	"pocok/src/db"
	"pocok/src/utils"
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

func (d *dependencies) handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	invoices, _ := db.GetPendingInvoices(d.dbClient, d.tableName)
	// invoices := mock.Invoices
	response := models.InvoiceResponse{
		Items: invoices,
		Total: len(invoices),
	}
	invoiceBytes, _ := json.Marshal(response)
	invoiceStr := string(invoiceBytes)
	return utils.ApiResponse(http.StatusOK, invoiceStr)
}

func main() {
	d := &dependencies{
		tableName: os.Getenv("tableName"),
		dbClient:  aws_clients.GetDbClient(),
	}

	lambda.Start(d.handler)
}
