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
	claims, parseTokenError := auth.ParseToken(token)
	if parseTokenError != nil {
		return utils.MailApiResponse(http.StatusUnauthorized, ""), parseTokenError
	}

	invoices, getPendingInvoicesError := db.GetPendingInvoices(d.dbClient, d.tableName, claims.OrgId)
	if getPendingInvoicesError != nil {
		utils.LogError("Error while getting pending invoices from db", getPendingInvoicesError)
		return nil, getPendingInvoicesError
	}
	// invoices := mocks.Invoices
	response := models.InvoiceResponse{
		Items: invoices,
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
