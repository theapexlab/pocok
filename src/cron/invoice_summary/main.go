package main

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"os"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type dependencies struct {
	tableName string
	queueUrl  string
	sqsClient *sqs.Client
	dbClient  *dynamodb.Client
}

// TODO refactor later to separate db file
func GetPendingInvoices(d *dependencies) ([]models.Invoice, error) {
	// TODO query the pending ones
	resp, err := d.dbClient.BatchGetItem(context.TODO(), &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			// "Invoices": types.KeysAndAttributes{},
		},
	})
	if err != nil {
		utils.LogError("Error while querying the db", err)
		return []models.Invoice{}, err
	}

	invoiceTable := resp.Responses["Invoices"]

	invoices := []models.Invoice{}
	for _, item := range invoiceTable {
		invoice := models.Invoice{}
		err := attributevalue.UnmarshalMap(item, invoice)
		if err != nil {
			utils.LogError("Error while loading invoices", err)
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}

func GetBody(invoices []models.Invoice) (string, error) {
	ids := make([]string, len(invoices))
	for i, inv := range invoices {
		ids[i] = inv.Id
	}
	var templateBuffer bytes.Buffer
	t, err := template.ParseFiles("src/utils/email_template.html")
	if err != nil {
		return "", err
	}
	execerr := t.Execute(&templateBuffer, ids)
	if execerr != nil {
		return "", err
	}
	return templateBuffer.String(), nil
}

func GetAttachments(invoices []models.Invoice) []string {
	attachments := []string{}
	for _, invoice := range invoices {
		attachments = append(attachments, invoice.Filename)
	}
	return attachments
}

func CreateEmail(invoices []models.Invoice) (*models.Email, error) {
	to := "billing@apexlab.io"
	subject := "Pocok Invoice Summary"
	html, err := GetBody(invoices)
	if err != nil {
		return nil, err
	}
	attachments := GetAttachments(invoices)

	email := models.Email{
		To:          to,
		Subject:     subject,
		Html:        html,
		Attachments: attachments,
	}
	return &email, nil
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
	invoices, invoiceErr := GetPendingInvoices(d)
	if invoiceErr != nil {
		utils.LogError("Error while loading invoices", invoiceErr)
		return invoiceErr
	}
	email, emailErr := CreateEmail(invoices)
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
