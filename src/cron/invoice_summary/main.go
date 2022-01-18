package main

import (
	"context"
	"encoding/json"
	"os"
	"pocok/src/cron/invoice_summary/create_email"
	"pocok/src/db"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type dependencies struct {
	tableName string
	queueUrl  string
	sqsClient *sqs.Client
	dbClient  *dynamodb.Client
}

func QueueEmail(d *dependencies, email models.Email) error {
	emailJson, jsonErr := json.Marshal(email)
	if jsonErr != nil {
		utils.LogError("Error while converting email struct to json", jsonErr)
		return jsonErr
	}
	message := string(emailJson)
	_, sqsErr := d.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: &message,
		QueueUrl:    &d.queueUrl,
	})
	if sqsErr != nil {
		utils.LogError("Error while sending message to SQS", sqsErr)
		return sqsErr
	}
	return nil
}

func (d *dependencies) handler(event events.CloudWatchEvent) error {
	invoices, invoiceErr := db.GetPendingInvoices(d.dbClient, d.tableName)
	if invoiceErr != nil {
		utils.LogError("Error while loading invoices", invoiceErr)
		return invoiceErr
	}
	email, emailErr := create_email.CreateEmail(invoices)
	if emailErr != nil {
		utils.LogError("Error while creating email", emailErr)
	}
	err := QueueEmail(d, *email)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	d := &dependencies{
		tableName: os.Getenv("tableName"),
		queueUrl:  os.Getenv("queueUrl"),
		dbClient:  aws_clients.GetDbClient(),
		sqsClient: aws_clients.GetSQSClient(),
	}

	lambda.Start(d.handler)
}
