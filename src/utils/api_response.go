package utils

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
)

func ApiResponse(status int, body string) *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{}
	resp.StatusCode = status
	resp.Body = body
	return &resp
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
