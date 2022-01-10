package main

import (
	"context"
	"net/http"
	"os"
	"pocok/src/api/process_email/get_pdf_url"
	"pocok/src/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type dependencies struct {
	queueUrl  string
	sqsClient *sqs.Client
}

func (d *dependencies) handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	pdfUrl, err := get_pdf_url.GetPdfUrl(request.Body)
	if err != nil {
		return utils.ApiResponse(http.StatusInternalServerError, "")
	}

	d.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: &pdfUrl,
		QueueUrl:    &d.queueUrl,
	})

	return utils.ApiResponse(http.StatusOK, "")
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
