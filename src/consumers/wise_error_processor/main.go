package main

import (
	"encoding/json"
	"os"
	"pocok/src/db"
	"pocok/src/services/wise"
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

func main() {
	d := &dependencies{
		dbClient:  aws_clients.GetDbClient(),
		tableName: os.Getenv("tableName"),
	}

	lambda.Start(d.handler)
}

func (d *dependencies) handler(event events.SQSEvent) error {
	for _, record := range event.Records {
		var messageData wise.WiseMessageData
		if unmarshalError := json.Unmarshal([]byte(record.Body), &messageData); unmarshalError != nil {
			utils.LogError("handler - Unmarshal", unmarshalError)
			return unmarshalError
		}

		updateError := db.UpdateInvoiceStatus(d.dbClient, d.tableName, models.APEX_ID, db.StatusUpdate{
			InvoiceId: messageData.Invoice.InvoiceId,
			Status:    models.TRANSFER_ERROR,
		})
		if updateError != nil {
			utils.LogError("error while updating invoice", updateError)
			return updateError
		}
	}
	return nil
}
