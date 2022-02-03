package main

import (
	"context"
	"os"
	"pocok/src/consumers/email_sender/create_email"
	"pocok/src/services/mailgun"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/aws_utils"
	"pocok/src/utils/models"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type dependencies struct {
	mailgunSender   string
	mailgunDomain   string
	mailgunApiKey   string
	emailRecipient  string
	apiUrl          string
	s3Client        *s3.Client
	assetBucketName string
}

func main() {
	d := &dependencies{
		mailgunSender:   os.Getenv("mailgunSender"),
		mailgunDomain:   os.Getenv("mailgunDomain"),
		mailgunApiKey:   os.Getenv("mailgunApiKey"),
		emailRecipient:  os.Getenv("emailRecipient"),
		apiUrl:          os.Getenv("apiUrl"),
		s3Client:        aws_clients.GetS3Client(),
		assetBucketName: os.Getenv("assetBucketName"),
	}

	lambda.Start(d.handler)
}

func (d *dependencies) handler(event events.SQSEvent) error {
	for _, record := range event.Records {
		if record.Body == models.EMAIL_SUMMARY {
			sendInvoiceError := SendInvoiceSummary(d)
			if sendInvoiceError != nil {
				return sendInvoiceError
			}
		}
	}

	return nil
}

func SendInvoiceSummary(d *dependencies) error {
	subject := models.EMAIL_SUMMARY_SUBJECT
	emailData, createEmailError := CreateEmail(d)
	if createEmailError != nil {
		utils.LogError("Error while creating email", createEmailError)
		return createEmailError
	}
	sendEmailError := SendEmail(d, subject, emailData)
	if sendEmailError != nil {
		utils.LogError("Error while sending email", sendEmailError)
		return createEmailError
	}
	return nil
}

func CreateEmail(d *dependencies) (string, error) {
	logoKey := "pocok-logo.png"
	logoUrl, getAssetUrlError := aws_utils.GetAssetUrl(*d.s3Client, d.assetBucketName, logoKey)
	if getAssetUrlError != nil {
		return "", getAssetUrlError
	}

	amp, getHtmlSummaryError := create_email.GetHtmlSummary(d.apiUrl, logoUrl)
	if getHtmlSummaryError != nil {
		return "", getHtmlSummaryError
	}

	return amp, nil
}

func SendEmail(d *dependencies, subject string, amp string) error {
	// Create an instance of the Mailgun Client
	client := mailgun.GetClient(d.mailgunDomain, d.mailgunApiKey)

	// The message object allows you to add attachments and Bcc recipients
	message := client.NewMessage(d.mailgunSender, subject, models.EMAIL_NO_AMP_BODY, d.emailRecipient)

	message.SetAMPHtml(amp)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	_, _, sendEmailError := client.Send(ctx, message)

	if sendEmailError != nil {
		utils.LogError("Failed to send email", sendEmailError)
		return sendEmailError
	}

	return nil
}
