package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type dependencies struct {
	queueUrl  string
	sqsClient *sqs.Client
}

func (d *dependencies) handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := d.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String("test message"),
		QueueUrl:    &d.queueUrl,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       *resp.MessageId,
	}, nil
}

func getSQSClient() *sqs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	return sqs.NewFromConfig(cfg)
}

func main() {
	d := &dependencies{
		queueUrl:  os.Getenv("queueUrl"),
		sqsClient: getSQSClient(),
	}

	lambda.Start(d.handler)
}
