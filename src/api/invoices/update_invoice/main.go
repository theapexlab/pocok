package main

import (
	"net/http"
	"os"
	"pocok/src/db"
	"pocok/src/utils"
	"pocok/src/utils/auth"
	"pocok/src/utils/aws_clients"

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
	jwt := r.QueryStringParameters["jwt"]
	claims, err := auth.ParseJwt(jwt)
	if err != nil {
		return utils.MailApiResponse(http.StatusUnauthorized, ""), err
	}

	id := ""
	status := ""
	updateErr := db.UpdateInvoiceStatus(d.dbClient, d.tableName, claims.OrgId, id, status)
	if updateErr != nil {
		return utils.MailApiResponse(http.StatusInternalServerError, ""), nil
	}

	return utils.MailApiResponse(http.StatusOK, ""), nil
}
