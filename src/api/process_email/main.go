package main

import (
	"context"
	"fmt"
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
	pdfUrl, pdfParseErr := get_pdf_url.GetPdfUrl(request.Body)
	if pdfParseErr != nil {
		fmt.Printf("❌ Error while parsing PDF URL: %s", pdfParseErr)
		return utils.ApiResponse(http.StatusInternalServerError, "")
	}

	_, sqsErr := d.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: &pdfUrl,
		QueueUrl:    &d.queueUrl,
	})
	if sqsErr != nil {
		fmt.Printf("❌ Error while sending message to SQS: %s", sqsErr)
		return utils.ApiResponse(http.StatusInternalServerError, "")
	}

	return utils.ApiResponse(http.StatusOK, "")
}

func main() {
	d := &dependencies{
		queueUrl:  os.Getenv("queueUrl"),
		sqsClient: aws_clients.GetSQSClient(),
	}

	lambda.Start(d.handler)
}
