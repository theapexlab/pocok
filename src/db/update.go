package db

import (
	"context"
	"errors"
	"pocok/src/utils"
	"pocok/src/utils/models"
	"strings"

	"github.com/almerlucke/go-iban/iban"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type StatusUpdate struct {
	InvoiceId string
	Status    string
}

func CreateValidStatusUpdate(data map[string]string) (StatusUpdate, error) {
	var update StatusUpdate
	err := utils.MapToStruct(data, &update)
	if err != nil {
		return update, errors.New("invalid input")
	}
	if update.InvoiceId == "" {
		return update, errors.New("invalid invoiceId")
	}
	if update.Status != models.ACCEPTED && update.Status != models.REJECTED {
		return update, errors.New("invalid status")
	}
	return update, nil
}

func UpdateInvoiceStatus(client *dynamodb.Client, tableName string, orgId string, update StatusUpdate) error {
	_, err := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			"sk": &types.AttributeValueMemberS{Value: models.INVOICE + "#" + update.InvoiceId},
		},
		UpdateExpression: aws.String("set #k1 = :v1, #k2 = :v2"),
		ExpressionAttributeNames: map[string]string{
			"#k1": "status",
			"#k2": "lsi1sk",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: update.Status},
			":v2": &types.AttributeValueMemberS{Value: models.STATUS + "#" + update.Status},
		},
	})
	return err
}

func CreateValidDataUpdate(data map[string]string) (models.Invoice, error) {
	var update models.Invoice
	err := utils.MapToStruct(data, &update)

	index := 0
	for {
		service := models.Service{}
		serviceMap := map[string]string{}
		found := false
		for key, val := range data {
			parts := strings.Split(key, "_")
			if strings.HasPrefix(parts[0], "service") && parts[2] == string(index) {
				found = true
				fieldName := parts[1]
				serviceMap[fieldName] = val
			}
		}
		if !found {
			break
		}
		err := utils.MapToStruct(service, &service)
		if err != nil {
			utils.LogError("error while parsing service", err)
		}
		update.Services = append(update.Services, service)
		index++
	}

	if err != nil {
		return update, errors.New("invalid input")
	}
	if update.InvoiceId == "" {
		return update, errors.New("invalid invoiceId")
	}

	if update.AccountNumber == "" && update.Iban == "" {
		// Todo: is it valid to provide accountNumber AND Iban
	}

	if update.Iban != "" {
		_, err := iban.NewIBAN(update.Iban)
		if err != nil {
			return update, errors.New("invalid iban")
		}
	}

	if update.AccountNumber != "" {
		_, err := iban.NewIBAN(update.Iban)
		if err != nil {
			return update, errors.New("invalid iban")
		}
	}

	//  vendorEmail => check if valid & update in db
	//  invoiceNumber =>  invoiceNumber + vendorName must be unique
	//  currency => eur, huf , usd
	//  dueDate =>  musst be in future
	//  grossAmount =>  >0,

	return update, nil
}

func UpdateInvoiceData(client *dynamodb.Client, tableName string, updateInvoice models.Invoice) error {
	_, err := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			"sk": &types.AttributeValueMemberS{Value: models.INVOICE + "#" + update.InvoiceId},
		},
		UpdateExpression: aws.String("set #k1 = :v1, #k2 = :v2"),
		ExpressionAttributeNames: map[string]string{
			"#k1": "status",
			"#k2": "lsi1sk",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: update.Status},
			":v2": &types.AttributeValueMemberS{Value: models.STATUS + "#" + update.Status},
		},
	})
	return err
}
