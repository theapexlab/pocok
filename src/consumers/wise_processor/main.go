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

		wiseDeps := wise.WiseDependencies{
			WiseService:  d.wiseService,
			SqsClient:    d.sqsClient,
			WiseQueueUrl: d.wiseQueueUrl,
		}

		step4Error := wiseDeps.Step4CreateTransfer(messageData)
		if step4Error != nil {
			utils.LogError("handler - step4", step4Error)
			return step4Error
		}

		updateError := db.UpdateInvoiceStatus(d.dbClient, d.tableName, db.UpdateStatusInput{
			OrgId:     models.APEX_ID,
			InvoiceId: messageData.Invoice.InvoiceId,
			Status:    models.ACCEPTED,
		})
		if updateError != nil {
			utils.LogError("error while updating invoice", updateError)
		}
	}
	return nil
}
