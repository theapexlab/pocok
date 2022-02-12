package main

import (
	"encoding/json"
	"os"
	"pocok/src/db"
	"pocok/src/services/wise"
	apiModels "pocok/src/services/wise/api/models"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type dependencies struct {
	wiseQueueUrl string
	sqsClient    *sqs.Client
	wiseService  *wise.WiseService
	dbClient     *dynamodb.Client
	tableName    string
}

func main() {
	d := &dependencies{
		wiseQueueUrl: os.Getenv("queueUrl"),
		sqsClient:    aws_clients.GetSQSClient(),
		wiseService:  wise.CreateWiseService(os.Getenv("wiseApiToken")),
		dbClient:     aws_clients.GetDbClient(),
		tableName:    os.Getenv("tableName"),
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

		step4Error := d.step4CreateTransfer(messageData)
		if step4Error != nil {
			utils.LogError("handler - step4", step4Error)
			return step4Error
		}
		// TODO dead letter queue

		db.UpdateInvoiceStatus(d.dbClient, d.tableName, models.APEX_ID, db.StatusUpdate{
			InvoiceId: messageData.Invoice.InvoiceId,
			Status:    models.ACCEPTED,
		})
	}
	return nil
}

func (d *dependencies) step4CreateTransfer(step4Data wise.WiseMessageData) error {
	transferInput := apiModels.Transfer{
		TargetAccount: step4Data.RecipientAccountId,
		QuoteUUID:     step4Data.QuoteId,
		Details: struct {
			Reference string `json:"reference"`
		}{Reference: step4Data.Invoice.InvoiceNumber},
		CustomerTransactionID: step4Data.TransactionId,
	}
	_, createTransferError := d.wiseService.WiseApi.CreateTransfer(transferInput)
	if createTransferError != nil {
		utils.LogError("step4CreateTransfer - CreateTransfer", createTransferError)
		return createTransferError
	}

	return nil
}
