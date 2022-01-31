package main

import (
	"errors"
	"net/http"
	"os"
	"pocok/src/utils"
	"pocok/src/utils/auth"
	"pocok/src/utils/aws_clients"
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

	validateError := validateUpdateDataRequest(data)
	if validateError != nil {
		utils.LogError("Invalid update payload", validateError)
		return utils.MailApiResponse(http.StatusUnprocessableEntity, validateError.Error()), nil
	}

}

func validateUpdateDataRequest(data map[string]string) error {
	invoiceId := data["invoiceId"]
	if invoiceId == "" {
		return errors.New("invalid invoiceId")
	}

	// status := data["status"]
	// if status != models.ACCEPTED && status != models.REJECTED {
	// 	return errors.New("invalid status")
	// }
	return nil
}
