package main

import (
	"context"
	"os"
	. "pocok/src/consumers/email_sender/create_email"
	"pocok/src/db"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/mailgun_client"
	"pocok/src/utils/models"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type dependencies struct {
	sender         string
	emailRecipient string
	apiUrl         string
	bucketName     string
	s3Client       *s3.Client
	tableName      string
	dbClient       *dynamodb.Client
}

func main() {
	d := &dependencies{
		sender:         os.Getenv("sender"),
		emailRecipient: os.Getenv("emailRecipient"),
		apiUrl:         os.Getenv("apiUrl"),
		bucketName:     os.Getenv("bucketName"),
		s3Client:       aws_clients.GetS3Client(),
		tableName:      os.Getenv("tableName"),
		dbClient:       aws_clients.GetDbClient(),
	}

	lambda.Start(d.handler)
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

func SendInvoiceSummary(d *dependencies) error {
	subject := models.EMAIL_SUMMARY_SUBJECT
	emailData, err := CreateEmail(d)
	if err != nil {
		utils.LogError("Error while creating email", err)
		return err
	}
	sendingErr := SendEmail(d, subject, emailData)
	if sendingErr != nil {
		utils.LogError("Error while sending email", sendingErr)
		return err
	}
	return nil
}

func CreateEmail(d *dependencies) (*models.EmailResponseData, error) {
	invoices, err := db.GetPendingInvoices(d.dbClient, d.tableName)
	if err != nil {
		utils.LogError("Error while loading invoices", err)
		return nil, err
	}

	amp, err := GetHtmlSummary(d.apiUrl)
	if err != nil {
		return nil, err
	}
	attachments, err := GetAttachments(d.s3Client, d.bucketName, invoices)
	if err != nil {
		return nil, err
	}

	response := models.EmailResponseData{
		Amp:         amp,
		Attachments: attachments,
	}
	return &response, nil
}

func SendEmail(d *dependencies, subject string, data *models.EmailResponseData) error {
	// Create an instance of the Mailgun Client
	client := mailgun_client.GetMailgunClient()

	// The message object allows you to add attachments and Bcc recipients
	message := client.NewMessage(d.sender, subject, models.EMAIL_NO_AMP_BODY, d.emailRecipient)
	for filename, attachment := range data.Attachments {
		message.AddBufferAttachment(filename, attachment)
	}
	message.SetAMPHtml(data.Amp)
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
