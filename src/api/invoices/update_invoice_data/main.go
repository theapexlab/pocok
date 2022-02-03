package main

import (
	"encoding/json"
	"net/http"
	"os"
	"pocok/src/db"
	"pocok/src/utils"
	"pocok/src/utils/auth"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"
	"pocok/src/utils/request_parser"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type dependencies struct {
	dbClient  *dynamodb.Client
	tableName string
}

func main() {
	d := &dependencies{
		tableName: os.Getenv("tableName"),
		dbClient:  aws_clients.GetDbClient(),
	}
	lambda.Start(d.handler)
}

func (d *dependencies) handler(r events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	token := r.QueryStringParameters["token"]
	claims, parseTokenError := auth.ParseToken(token)
	if parseTokenError != nil {
		utils.LogError("Token validation failed", parseTokenError)
		return utils.MailApiResponse(http.StatusUnauthorized, ""), parseTokenError
	}

	data, parseFormDataError := request_parser.ParseUrlEncodedFormData(r)
	if parseFormDataError != nil {
		utils.LogError("Form body parse failed", parseFormDataError)
		return utils.MailApiResponse(http.StatusBadRequest, ""), parseFormDataError
	}

	update, validationErr := db.CreateValidDataUpdate(data)
	if validationErr != nil {
		utils.LogError("Invalid update payload", validationErr)
		response := models.ValidationErrorResponse{
			Message: validationErr.Error(),
		}
		messageBytes, marshalError := json.Marshal(response)
		if marshalError != nil {
			utils.LogError("Error while parsing validation error response", marshalError)
			return nil, marshalError
		}

		messageStr := string(messageBytes)
		return utils.MailApiResponse(http.StatusUnprocessableEntity, messageStr), nil
	}

	updateErr := db.UpdateInvoiceData(d.dbClient, d.tableName, claims.OrgId, update)
	if updateErr != nil {
		utils.LogError("Error updating dynamo db", updateErr)
		return utils.MailApiResponse(http.StatusInternalServerError, ""), nil
	}

	return utils.MailApiResponse(http.StatusOK, ""), nil

}
