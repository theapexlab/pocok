package main

import (
	"net/http"
	"os"
	"pocok/src/api/invoices/update_utils"
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
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type dependencies struct {
	dbClient       *dynamodb.Client
	tableName      string
	s3Client       *s3.Client
	bucketName     string
	typlessToken   string
	typlessDocType string
	wiseQueueUrl   string
	sqsClient      *sqs.Client
}

func main() {
	d := &dependencies{
		dbClient:       aws_clients.GetDbClient(),
		tableName:      os.Getenv("tableName"),
		s3Client:       aws_clients.GetS3Client(),
		bucketName:     os.Getenv("bucketName"),
		typlessToken:   os.Getenv("typlessToken"),
		typlessDocType: os.Getenv("typlessDocType"),
		wiseQueueUrl:   os.Getenv("wiseQueueUrl"),
		sqsClient:      aws_clients.GetSQSClient(),
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

	statusUpdate, validationError := db.GetValidStatusUpdate(data)
	if validationError != nil {
		utils.LogError("Invalid while validating update", validationError)
		return utils.MailApiResponse(http.StatusUnprocessableEntity, validationError.Error()), nil
	}

	if statusUpdate.Status == models.REJECTED {
		return d.rejectInvoice(*claims, *statusUpdate)
	}
	if statusUpdate.Status == models.ACCEPTED {
		return d.acceptInvoice(*claims, *statusUpdate)
	}

	return utils.MailApiResponse(http.StatusTeapot, ""), nil
}

func (d *dependencies) rejectInvoice(claims models.JWTClaims, update db.StatusUpdate) (*events.APIGatewayProxyResponse, error) {
	deleteError := db.DeleteInvoice(d.dbClient, d.tableName, *d.s3Client, d.bucketName, claims.OrgId, update.InvoiceId, update.Filename)
	if deleteError != nil {
		utils.LogError("Error while removing invoice", deleteError)
		return utils.MailApiResponse(http.StatusInternalServerError, ""), nil
	}
	return utils.MailApiResponse(http.StatusOK, ""), nil
}

func (d *dependencies) acceptInvoice(claims models.JWTClaims, update db.StatusUpdate) (*events.APIGatewayProxyResponse, error) {
	updateError := db.UpdateInvoiceStatus(d.dbClient, d.tableName, claims.OrgId, update)
	if updateError != nil {
		utils.LogError("Error while updating invoice", updateError)
		return utils.MailApiResponse(http.StatusInternalServerError, ""), nil
	}
	invoice, getInvoiceError := db.GetInvoice(d.dbClient, d.tableName, claims.OrgId, update.InvoiceId)
	if getInvoiceError != nil {
		utils.LogError("Error while getting invoice", getInvoiceError)
		return utils.MailApiResponse(http.StatusOK, ""), nil
	}

	feedbackError := update_utils.UpdateTypeless(d.typlessToken, d.typlessDocType, *invoice)
	if feedbackError != nil {
		utils.LogError("Error while submitting typless feedback", feedbackError)
	}

	wiseError := update_utils.SendWiseMessage(*d.sqsClient, d.wiseQueueUrl, *invoice)
	if wiseError != nil {
		utils.LogError("Error while creating wise request", wiseError)
	}
	return utils.MailApiResponse(http.StatusOK, ""), nil
}
