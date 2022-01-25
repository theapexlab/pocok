package main

import (
	"errors"
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
	claims, err := auth.ParseToken(token)
	if err != nil {
		utils.LogError("Token validation failed", err)
		return utils.MailApiResponse(http.StatusUnauthorized, ""), err
	}

	data, err := request_parser.ParseUrlEncodedFormData(r)
	if err != nil {
		utils.LogError("Form body parse failed", err)
		return utils.MailApiResponse(http.StatusBadRequest, ""), err
	}

	validateError := validateUpdateRequest(data)
	if validateError != nil {
		utils.LogError("Invalid update payload", validateError)
		return utils.MailApiResponse(http.StatusUnprocessableEntity, validateError.Error()), nil
	}

	updateErr := db.UpdateInvoiceStatus(d.dbClient, d.tableName, claims.OrgId, data["invoiceId"], data["status"])
	if updateErr != nil {
		utils.LogError("Error updating dynamo db", updateErr)
		return utils.MailApiResponse(http.StatusInternalServerError, ""), nil
	}
	return utils.MailApiResponse(http.StatusOK, ""), nil
}

func validateUpdateRequest(data map[string]string) error {
	invoiceId := data["invoiceId"]
	status := data["status"]
	if invoiceId == "" {
		return errors.New("invalid invoiceId")
	}
	if status != models.ACCEPTED && status != models.REJECTED {
		return errors.New("invalid status")
	}
	return nil
}