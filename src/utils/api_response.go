package utils

import (
	"github.com/aws/aws-lambda-go/events"
)

func ApiResponse(status int, body string) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{}
	resp.StatusCode = status
	resp.Body = body

	return &resp, nil
}
