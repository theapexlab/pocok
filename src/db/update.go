package db

import (
	"context"
	"encoding/json"
	"errors"
	"pocok/src/utils"
	"pocok/src/utils/models"
	"strings"

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
		// todo: is it valid to provide accountNumber AND Iban?
	}

	if update.Iban != "" {
		_, err := utils.ValidateIban(update.Iban)
		if err != nil {
			return update, errors.New("invalid iban")
		}
	}

	if update.AccountNumber != "" {
		_, err := utils.ValidateAccountNumber(update.AccountNumber)
		if err != nil {
			return update, err
		}
	}

	//  todo: add additional validation for fields below
	//  vendorEmail => check if valid & update in db
	//  invoiceNumber =>  invoiceNumber + vendorName must be unique
	//  currency => eur, huf , usd
	//  dueDate =>  musst be in future
	//  grossAmount =>  >0,

	return update, nil
}

func UpdateInvoiceData(client *dynamodb.Client, tableName string, orgId string, update models.Invoice) error {
	servicesJson, marshalErr := json.Marshal(update.Services)
	if marshalErr != nil {
		utils.LogError("Cannot marshal invoice services", marshalErr)
		return marshalErr
	}
	_, err := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			"sk": &types.AttributeValueMemberS{Value: models.INVOICE + "#" + update.InvoiceId},
		},
		UpdateExpression: aws.String(`set #k1 = :v1, #k2 = :v2, #k3 = :v3, #k4 = :v4, #k5 = :v5,
		 #k6 = :v6, #k7 = :v7, #k8 = :v8, #k9 = :v9, #k10 = :v10, #k11 = :v11`),
		ExpressionAttributeNames: map[string]string{
			"#k1":  "vendorName",
			"#k2":  "accountNumber",
			"#k3":  "iban",
			"#k4":  "netPrice",
			"#k5":  "grossPrice",
			"#k6":  "vatAmount",
			"#k7":  "vatRate",
			"#k8":  "currency",
			"#k9":  "dueDate",
			"#k10": "services",
			"#k11": "vendorEmail",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1":  &types.AttributeValueMemberS{Value: update.VendorName},
			":v2":  &types.AttributeValueMemberS{Value: update.AccountNumber},
			":v3":  &types.AttributeValueMemberS{Value: update.Iban},
			":v4":  &types.AttributeValueMemberS{Value: update.NetPrice},
			":v5":  &types.AttributeValueMemberS{Value: update.GrossPrice},
			":v6":  &types.AttributeValueMemberS{Value: update.VatAmount},
			":v7":  &types.AttributeValueMemberS{Value: update.VatRate},
			":v8":  &types.AttributeValueMemberS{Value: update.Currency},
			":v9":  &types.AttributeValueMemberS{Value: update.DueDate},
			":v10": &types.AttributeValueMemberS{Value: string(servicesJson)},
			":v11": &types.AttributeValueMemberS{Value: update.VendorEmail},
		},
	})
	return err
}
