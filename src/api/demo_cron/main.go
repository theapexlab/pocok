package main

import (
	"context"
	"errors"
	"os"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type dependencies struct {
	queueUrl  string
	sqsClient *sqs.Client
	demoToken string
}

func (d *dependencies) handler(r events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	token := r.QueryStringParameters["token"]
	if token != d.demoToken {
		utils.LogError("Token validation failed", errors.New("token validation failed"))
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Token validation failed",
		}, nil
	}

	_, sqsErr := d.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(models.EMAIL_SUMMARY),
		QueueUrl:    &d.queueUrl,
	})
	if sqsErr != nil {
		utils.LogError("Error while sending message to SQS", sqsErr)
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       sqsErr.Error(),
		}, sqsErr
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Done",
	}, nil
}

func main() {
	d := &dependencies{
		queueUrl:  os.Getenv("queueUrl"),
		sqsClient: aws_clients.GetSQSClient(),
		demoToken: os.Getenv("demoToken"),
	}

	lambda.Start(d.handler)
}
