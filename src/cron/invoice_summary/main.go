package main

import (
	"context"
	"os"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type dependencies struct {
	queueUrl  string
	sqsClient *sqs.Client
}

func QueueEmailEvent(d *dependencies, emailType string) error {
	_, sqsErr := d.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: &emailType,
		QueueUrl:    &d.queueUrl,
	})
	if sqsErr != nil {
		utils.LogError("Error while sending message to SQS", sqsErr)
		return sqsErr
	}
	return nil
}

func (d *dependencies) handler(event events.CloudWatchEvent) error {
	err := QueueEmailEvent(d, models.EMAIL_SUMMARY)
	if err != nil {

		return err
	}
	return nil
}

func main() {
	d := &dependencies{
		queueUrl:  os.Getenv("queueUrl"),
		sqsClient: aws_clients.GetSQSClient(),
	}

	lambda.Start(d.handler)
}
