package main

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"os"
	"pocok/src/db"
	"pocok/src/utils"
	"pocok/src/utils/auth"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"

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
	data, err := parseUrlEncodedFormData(r)
	if err != nil {
		utils.LogError("Form body parse failed", err)
		return utils.MailApiResponse(http.StatusBadRequest, ""), err
	}
	invoiceId := data["invoiceId"]
	status := data["status"]
	if invoiceId == "" || (status != models.ACCEPTED && status != models.REJECTED) {
		utils.LogError("Invalid update payload", errors.New(""))
		return utils.MailApiResponse(http.StatusUnprocessableEntity, ""), nil
	}
	updateErr := db.UpdateInvoiceStatus(d.dbClient, d.tableName, claims.OrgId, invoiceId, status)
	if updateErr != nil {
		utils.LogError("Error updating dynamo db", updateErr)
		return utils.MailApiResponse(http.StatusInternalServerError, ""), nil
	}
	return utils.MailApiResponse(http.StatusOK, ""), nil
}

func parseUrlEncodedFormData(r events.APIGatewayProxyRequest) (map[string]string, error) {
	data, err := base64.StdEncoding.DecodeString(r.Body)
	result := map[string]string{}
	if err != nil {
		utils.LogError("Error decoding base64 message content", err)
		return result, err
	}

	parsedFormData, parseErr := url.ParseQuery(string(data))
	if parseErr != nil {
		utils.LogError("Error parsing url encoded form data", parseErr)
		return result, parseErr
	}

	for key, value := range parsedFormData {
		result[key] = value[0]
	}

	return result, nil
}

// Turns out that we do not need this...
//
// func parseMultipartFormData(r events.APIGatewayProxyRequest) (map[string]string, error) {
// 	data, err := base64.StdEncoding.DecodeString(r.Body)
// 	result := map[string]string{}
// 	if err != nil {
// 		utils.LogError("Error decoding base64 message content", err)
// 		return result, err
// 	}
// 	reader := bytes.NewReader(data)
// 	contentType := r.Headers["content-type"]
// 	_, params, err := mime.ParseMediaType(contentType)
// 	if err != nil {
// 		utils.LogError("Error during parse form media type", err)
// 		return result, err
// 	}
// 	mr := multipart.NewReader(reader, params["boundary"])
// 	for {
// 		part, err := mr.NextPart()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return result, err
// 		}
// 		slurp, err := io.ReadAll(part)
// 		if err != nil {
// 			return result, err
// 		}
// 		key := part.FormName()
// 		value := string(slurp)
// 		result[key] = value
// 	}
// 	return result, nil
// }
