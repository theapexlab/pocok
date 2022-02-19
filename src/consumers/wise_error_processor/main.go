package main

import (
	"encoding/json"
	"os"
	"pocok/src/db"
	"pocok/src/services/slack_service"
	"pocok/src/services/wise"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type dependencies struct {
	dbClient    *dynamodb.Client
	slackClient slack_service.SlackClient
	tableName   string
}

func main() {
	d := &dependencies{
		dbClient: aws_clients.GetDbClient(),
		slackClient: slack_service.SlackClient{
			Url:      os.Getenv("slackWebhookUrl"),
			Username: os.Getenv("slackUsername"),
			Channel:  os.Getenv("slackChannel"),
		},
		tableName: os.Getenv("tableName"),
	}

	lambda.Start(d.handler)
}

func (d *dependencies) handler(event events.SQSEvent) error {
	slackMessage := "Wise transfer error.\n"
	for _, record := range event.Records {
		var messageData wise.WiseMessageData
		if unmarshalError := json.Unmarshal([]byte(record.Body), &messageData); unmarshalError != nil {
			utils.LogError("handler - Unmarshal", unmarshalError)
			slackMessage += unmarshalError.Error()
			continue
		}

		slackMessage += "InvoiceId: " + messageData.Invoice.InvoiceId + "\n"

		updateError := db.UpdateInvoiceStatus(d.dbClient, d.tableName, db.UpdateStatusInput{
			OrgId:     models.APEX_ID,
			InvoiceId: messageData.Invoice.InvoiceId,
			Status:    models.TRANSFER_ERROR,
		})
		if updateError != nil {
			utils.LogError("error while updating invoice", updateError)
			slackMessage += updateError.Error()
		}
	}
	slackMessage += "\nFor further information, see the logs in AWS"

	_, slackError := d.slackClient.SendMessage(slackMessage)
	if slackError != nil {
		utils.LogError("error while seding slack message", slackError)
	}
	return nil
}
