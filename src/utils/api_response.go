package utils

import (
	"github.com/aws/aws-lambda-go/events"
)

func ApiResponse(status int, body string) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{}
	resp.StatusCode = status
	resp.Body = body
	resp.Headers = map[string]string{
		"Content-Type":           "application/json",
		"AMP-Email-Allow-Sender": "billing@apexlab.io",
	}

	return &resp, nil
}
