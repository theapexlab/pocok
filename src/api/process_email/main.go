package main

import (
	"context"
	"net/http"
	"os"
	"pocok/src/api/process_email/get_pdf_url"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func main() {
	d := &dependencies{
		queueUrl:  os.Getenv("queueUrl"),
		sqsClient: aws_clients.GetSQSClient(),
	}

	lambda.Start(d.handler)
}
