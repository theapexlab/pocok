package main

import (
	"context"
	"os"
	"pocok/src/consumers/email_sender/create_email"
	"pocok/src/db"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mailgun/mailgun-go/v4"
)

type dependencies struct {
	domain     string
	sender     string
	apiKey     string
	bucketName string
	s3Client   *s3.Client
	tableName  string
	dbClient   *dynamodb.Client
}

func SendEmail(d *dependencies, recipient string, subject string, body string, attachments map[string][]byte) error {
	// Create an instance of the Mailgun Client
	client := mailgun.NewMailgun(d.domain, d.apiKey)

	// The message object allows you to add attachments and Bcc recipients
	message := client.NewMessage(d.sender, subject, "Yeet boi", recipient)
	for filename, attachment := range attachments {
		message.AddBufferAttachment(filename, attachment)
	}
	message.SetAMPHtml(body)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	_, _, sendErr := client.Send(ctx, message)

	if sendErr != nil {
		utils.LogError("Failed to send email", sendErr)
		return sendErr
	}

	return nil
}

func SendInvoiceSummary(d *dependencies) error {
	invoices, invoiceErr := db.GetPendingInvoices(d.dbClient, d.tableName)
	if invoiceErr != nil {
		utils.LogError("Error while loading invoices", invoiceErr)
		return invoiceErr
	}

	html, htmlErr := create_email.GetHtmlSummary(invoices)
	if htmlErr != nil {
		return htmlErr
	}
	attachments, attachmentErr := create_email.GetAttachments(d.s3Client, d.bucketName, invoices)
	if attachmentErr != nil {
		return attachmentErr
	}

	subject := "SUBJECT"
	recipient := "tom@apexlab.io"

	emailErr := SendEmail(d, recipient, subject, html, attachments)
	if emailErr != nil {
		utils.LogError("Error while sending email", emailErr)
		return emailErr
	}
	return nil
}

func (d *dependencies) handler(event events.SQSEvent) error {
	for _, record := range event.Records {
		if record.Body == models.EMAIL_SUMMARY {
			err := SendInvoiceSummary(d)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	d := &dependencies{
		domain:     os.Getenv("domain"),
		apiKey:     os.Getenv("apiKey"),
		sender:     os.Getenv("sender"),
		bucketName: os.Getenv("bucketName"),
		s3Client:   aws_clients.GetS3Client(),
		tableName:  os.Getenv("tableName"),
		dbClient:   aws_clients.GetDbClient(),
	}

	lambda.Start(d.handler)
}
