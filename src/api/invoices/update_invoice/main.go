package main

import (
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
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type dependencies struct {
	dbClient   *dynamodb.Client
	tableName  string
	s3Client   *s3.Client
	bucketName string
}

func main() {
	d := &dependencies{
		tableName:  os.Getenv("tableName"),
		dbClient:   aws_clients.GetDbClient(),
		bucketName: os.Getenv("bucketName"),
		s3Client:   aws_clients.GetS3Client(),
	}
	lambda.Start(d.handler)
}

// TODO? hidden _method param for delete & update requests
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

	update, validationErr := db.CreateStatusUpdate(data)
	if validationErr != nil {
		utils.LogError("Invalid update payload", validationErr)
		return utils.MailApiResponse(http.StatusUnprocessableEntity, validationErr.Error()), nil
	}

	if update.Status == models.REJECTED {
		removeErr := db.DeleteInvoice(d.dbClient, d.tableName, *d.s3Client, d.bucketName, claims.OrgId, update.InvoiceId)
		if removeErr != nil {
			utils.LogError("Error updating dynamo db", removeErr)
			return utils.MailApiResponse(http.StatusInternalServerError, ""), nil
		}
		return utils.MailApiResponse(http.StatusOK, ""), nil
	}

	if update.Status == models.ACCEPTED {
		updateErr := db.UpdateInvoiceStatus(d.dbClient, d.tableName, claims.OrgId, update)
		if updateErr != nil {
			utils.LogError("Error updating dynamo db", updateErr)
			return utils.MailApiResponse(http.StatusInternalServerError, ""), nil
		}
		return utils.MailApiResponse(http.StatusOK, ""), nil
	}

	return utils.MailApiResponse(http.StatusTeapot, ""), nil
}
