package utils

import (
	"encoding/json"
	"os"
	"pocok/src/utils"
	"pocok/src/utils/models"

	"github.com/aws/aws-lambda-go/events"
)

func ApiResponse(status int, body string) *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{}
	resp.StatusCode = status
	resp.Body = body
	return &resp
}

func ApiErrorBody(message string) string {
	response := models.ErrorResponse{
		Message: message,
	}
	jsonBytes, jsonError := json.Marshal(response)
	if jsonError != nil {
		utils.LogError("error while marshaling message", jsonError)
		return ""
	}
	return string(jsonBytes)
}

func MailApiResponse(status int, body string) *events.APIGatewayProxyResponse {
	resp := ApiResponse(status, body)
	if resp.Body == "" {
		resp.Body = "{}"
	}
	resp.Headers = map[string]string{
		"Content-Type":           "application/json",
		"AMP-Email-Allow-Sender": os.Getenv("mailgunSender"),
	}
	return resp
}
