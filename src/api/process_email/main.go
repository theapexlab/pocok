package main

import (
	"context"
	"encoding/json"
	"errors"
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
	invoiceMessages, emailParseError := parse_email.ParseEmail(request.Body)
	if emailParseError != nil {
		utils.LogError("error while parsing email", emailParseError)
		return utils.ApiResponse(http.StatusInternalServerError, ""), emailParseError
	}

	invoiceErrors := ""
	for _, invoiceMessage := range invoiceMessages {
		invoiceMessageByteArr, jsonError := json.Marshal(invoiceMessage)
		if jsonError != nil {
			utils.LogError("error while parsing json", jsonError)
			invoiceErrors += jsonError.Error()
			continue
		}
		invoiceMessageString := string(invoiceMessageByteArr)

		_, sqsError := d.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
			MessageBody: &invoiceMessageString,
			QueueUrl:    &d.queueUrl,
		})
		if sqsError != nil {
			utils.LogError("error while sending message to SQS", sqsError)
			invoiceErrors += sqsError.Error()
		}
	}

	if len(invoiceErrors) != 0 {
		requestError := errors.New(invoiceErrors)
		utils.LogError("error while handling request", requestError)
		return utils.ApiResponse(http.StatusInternalServerError, ""), requestError
	}

	return utils.ApiResponse(http.StatusOK, ""), nil
}

func main() {
	d := &dependencies{
		queueUrl:  os.Getenv("queueUrl"),
		sqsClient: aws_clients.GetSQSClient(),
	}

	lambda.Start(d.handler)
}
