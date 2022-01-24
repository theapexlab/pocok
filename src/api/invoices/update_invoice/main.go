package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
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
		return utils.MailApiResponse(http.StatusUnauthorized, ""), err
	}
	data, err := parseFormBody(r)
	if err != nil {
		return utils.MailApiResponse(http.StatusBadRequest, ""), nil
	}
	id := data["id"]
	status := data["status"]
	if id == "" || (status != models.ACCEPTED && status != models.REJECTED) {
		return utils.MailApiResponse(http.StatusUnprocessableEntity, ""), nil
	}
	updateErr := db.UpdateInvoiceStatus(d.dbClient, d.tableName, claims.OrgId, id, status)
	if updateErr != nil {
		utils.LogError("Error updating dynamo db", updateErr)
		return utils.MailApiResponse(http.StatusInternalServerError, ""), nil
	}
	return utils.MailApiResponse(http.StatusOK, ""), nil
}

func parseFormBody(r events.APIGatewayProxyRequest) (map[string]string, error) {
	data, err := base64.StdEncoding.DecodeString(r.Body)
	result := map[string]string{}
	if err != nil {
		utils.LogError("Error decoding base64 message content", err)
		return result, err
	}
	reader := bytes.NewReader(data)
	contentType := r.Headers["content-type"]
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		utils.LogError("Error during parse form media type", err)
		return result, err
	}
	mr := multipart.NewReader(reader, params["boundary"])
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return result, err
		}
		slurp, err := io.ReadAll(part)
		if err != nil {
			return result, err
		}
		key := part.FormName()
		value := string(slurp)
		result[key] = value
	}
	return result, nil
}
