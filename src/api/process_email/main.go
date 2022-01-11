package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"pocok/src/api/process_email/parse_email"
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
	invoiceMessage, emailParseErr := parse_email.ParseEmail(request.Body)
	if emailParseErr != nil {
		utils.LogError("Error while parsing email", emailParseErr)
		return utils.ApiResponse(http.StatusInternalServerError, "")
	}

	invoiceMessageByteArr, _ := json.Marshal(invoiceMessage)
	invoiceMessageString := string(invoiceMessageByteArr)

	_, sqsErr := d.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: &invoiceMessageString,
		QueueUrl:    &d.queueUrl,
	})
	if sqsErr != nil {
		utils.LogError("Error while sending message to SQS", emailParseErr)
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
