package update_utils

import (
	"context"
	"encoding/json"
	"pocok/src/services/wise"
	"pocok/src/utils"
	"pocok/src/utils/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func SendWiseMessage(sqsClient sqs.Client, wiseQueueUrl string, invoice models.Invoice) error {
	messageBody := wise.WiseMessageData{
		RequestType: wise.WiseStep1,
		Invoice:     invoice,
	}
	messageByteArray, marshalError := json.Marshal(messageBody)
	if marshalError != nil {
		utils.LogError("sendMessage - Marshal", marshalError)
		return marshalError
	}
	messageString := string(messageByteArray)
	_, sqsError := sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(wiseQueueUrl),
		MessageBody: aws.String(messageString),
	})
	if sqsError != nil {
		utils.LogError("error pushing message to queue", sqsError)
		return sqsError
	}
	return nil
}
