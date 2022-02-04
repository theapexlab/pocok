package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"pocok/src/db"
	"pocok/src/services/typless"
	"pocok/src/services/typless/create_training_data"
	"pocok/src/services/wise"
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
		deleteError := db.DeleteInvoice(d.dbClient, d.tableName, *d.s3Client, d.bucketName, claims.OrgId, statusUpdate.InvoiceId, statusUpdate.Filename)
		if deleteError != nil {
			utils.LogError("Error while removing invoice", deleteError)
			return utils.MailApiResponse(http.StatusInternalServerError, ""), nil
		}
		return utils.MailApiResponse(http.StatusOK, ""), nil
	}

	if statusUpdate.Status == models.ACCEPTED {
		updateError := db.UpdateInvoiceStatus(d.dbClient, d.tableName, claims.OrgId, *statusUpdate)
		if updateError != nil {
			utils.LogError("Error while updating invoice", updateError)
			return utils.MailApiResponse(http.StatusInternalServerError, ""), nil
		}
		feedbackError := updateFeedback(d, claims.OrgId, statusUpdate.InvoiceId)
		if feedbackError != nil {
			utils.LogError("Error while submitting typless feedback", feedbackError)
		}
		invoice, getInvoiceError := db.GetInvoice(d.dbClient, d.tableName, claims.OrgId, statusUpdate.InvoiceId)
		if getInvoiceError != nil {
			utils.LogError("Error while getting invoice", getInvoiceError)
		}
		messageBody := wise.WiseMessageData{
			RequestType: wise.WiseStep1,
			Invoice:     *invoice,
		}
		messageByteArray, marshalError := json.Marshal(messageBody)
		if marshalError != nil {
			utils.LogError("sendMessage - Marshal", marshalError)
			return nil, marshalError
		}
		messageString := string(messageByteArray)
		d.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
			QueueUrl:    &d.wiseQueueUrl,
			MessageBody: &messageString,
		})
		return utils.MailApiResponse(http.StatusOK, ""), nil
	}

	return utils.MailApiResponse(http.StatusTeapot, ""), nil
}

func updateFeedback(d *dependencies, orgId string, invoiceId string) error {
	invoice, getInvoiceError := db.GetInvoice(d.dbClient, d.tableName, orgId, invoiceId)
	if getInvoiceError != nil {
		utils.LogError("Error getting invoice from db", getInvoiceError)
		return getInvoiceError
	}

	typlessError := typless.AddDocumentFeedback(
		&typless.Config{
			Token:   d.typlessToken,
			DocType: d.typlessDocType,
		},
		*create_training_data.CreateTrainingData(invoice),
	)
	if typlessError != nil {
		utils.LogError("Error adding document feedback to typless", typlessError)
	}

	return nil
}
