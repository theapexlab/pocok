package main

import (
	"context"
	"encoding/json"
	"os"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mailgun/mailgun-go/v4"
)

type dependencies struct {
	domain     string
	apiKey     string
	sender     string
	recipient  string
	bucketName string
	s3Client   *s3.Client
}

func SendEmail(d *dependencies, subject string, body string, attachments map[string][]byte) error {
	// Create an instance of the Mailgun Client
	client := mailgun.NewMailgun(d.domain, d.apiKey)

	// The message object allows you to add attachments and Bcc recipients
	message := client.NewMessage(d.sender, subject, body, d.recipient)
	for filename, attachment := range attachments {
		message.AddBufferAttachment(filename, attachment)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	_, _, err := client.Send(ctx, message)

	if err != nil {
		utils.LogError("Failed to send email", err)
	}

	return err
}

func parseBody(body string) (*models.Email, error) {
	var jsonBody *models.Email

	if err := json.Unmarshal([]byte(body), &jsonBody); err != nil {
		return nil, err
	}

	return jsonBody, nil
}

func (d *dependencies) handler(event events.SQSEvent) error {
	for _, record := range event.Records {
		emailStruct, err := parseBody(record.Body)
		if err != nil {
			continue
		}
		attachments := map[string][]byte{}
		for _, attachmentString := range emailStruct.Attachments {
			_, s3Err := d.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
				Bucket: &d.bucketName,
				Key:    &attachmentString,
			})
			if s3Err != nil {
				continue
			}
			file := []byte{}
			if err != nil {
				continue
			}
			attachments[attachmentString] = file
		}
		emailErr := SendEmail(d, emailStruct.Subject, emailStruct.Html, attachments)
		if emailErr != nil {
			return emailErr
		}
	}

	return nil
}

func main() {
	d := &dependencies{
		domain:     os.Getenv("domain"),
		apiKey:     os.Getenv("apiKey"),
		sender:     os.Getenv("sender"),
		recipient:  os.Getenv("recipient"),
		bucketName: os.Getenv("bucketName"),
		s3Client:   aws_clients.GetS3Client(),
	}

	lambda.Start(d.handler)
}
