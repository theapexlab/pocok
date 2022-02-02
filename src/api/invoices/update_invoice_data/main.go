package main

import (
	"net/http"
	"os"
	"pocok/src/db"
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

	update, validationErr := db.CreateValidDataUpdate(data)
	if validationErr != nil {
		utils.LogError("Invalid update payload", validationErr)
		return utils.MailApiResponse(http.StatusUnprocessableEntity, validationErr.Error()), nil
	}

	updateErr := db.UpdateInvoiceData(d.dbClient, d.tableName, claims.OrgId, update)
	if updateErr != nil {
		utils.LogError("Error updating dynamo db", updateErr)
		return utils.MailApiResponse(http.StatusInternalServerError, ""), nil
	}

	//  todo: send json response
	// invoiceBytes, err := json.Marshal(response)
	// if err != nil {
	// 	utils.LogError("Error while parsing invoices from db", err)
	// 	return nil, err
	// }

	// invoiceStr := string(invoiceBytes)

	return utils.MailApiResponse(http.StatusOK, ""), nil

}
