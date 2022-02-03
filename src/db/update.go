package db

import (
	"context"
	"errors"
	"pocok/src/utils"
	"pocok/src/utils/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type StatusUpdate struct {
	InvoiceId string
	Status    string
	Filename  string
}

const TRANSACTION_LIMIT = 25 // 25 is the max number of items that can be updated in a single transaction

type VendorUpdate struct {
	VendorName  string
	VendorEmail string
}

func CreateStatusUpdate(data map[string]string) (*StatusUpdate, error) {
	var update StatusUpdate
	mapError := utils.MapToStruct(data, &update)
	if mapError != nil {
		return nil, errors.New("invalid input")
	}
	if update.InvoiceId == "" {
		return nil, errors.New("invalid invoiceId")
	}
	if update.Status != models.ACCEPTED && update.Status != models.REJECTED {
		return nil, errors.New("invalid status")
	}
	if update.Status == models.REJECTED && update.Filename == "" {
		return nil, errors.New("invalid update, filename must be present on reject")
	}
	return &update, nil
}

func UpdateInvoiceStatus(client *dynamodb.Client, tableName string, orgId string, update StatusUpdate) error {
	_, updateError := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
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
	return updateError
}

func CreateValidDataUpdate(data map[string]string) (models.Invoice, error) {
	update, updateError := utils.MapUpdateDataToInvoice(data)

	if updateError != nil {
		return update, errors.New("invalid input")
	}

	if update.InvoiceId == "" {
		return update, errors.New("invalid invoiceId")
	}

	if update.VendorName == "" {
		return update, errors.New("invalid vendor name")
	}

	if update.AccountNumber == "" && update.Iban == "" {
		return update, errors.New("iban or account number must be provided")
	}

	if update.Iban != "" {
		_, ibanError := utils.GetValidIban(update.Iban)
		if ibanError != nil {
			return update, ibanError
		}
	}

	if update.AccountNumber != "" {
		_, accountNumberError := utils.GetValidAccountNumber(update.AccountNumber)
		if accountNumberError != nil {
			return update, accountNumberError
		}
	}

	_, priceError := utils.GetValidPrice(update.GrossPrice)
	if priceError != nil {
		return update, errors.New("invalid gross price")
	}

	_, dateError := utils.GetValidDueDate(update.DueDate)
	if dateError != nil {
		return update, errors.New("invalid due date")
	}

	_, currError := utils.GetValidCurrency(update.Currency)
	if currError != nil {
		return update, currError
	}

	//  todo: add additional validation for fields below
	//  vendorEmail => update vendor email if changed
	//  invoiceNumber =>  invoiceNumber + vendorName must be unique

	return update, nil
}

func UpdateInvoiceData(client *dynamodb.Client, tableName string, orgId string, update models.Invoice) error {
	serviceList, marshalError := attributevalue.MarshalList(update.Services)
	if marshalError != nil {
		utils.LogError("Cannot marshal invoice services", marshalError)
		return marshalError
	}

	_, updateError := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			"sk": &types.AttributeValueMemberS{Value: models.INVOICE + "#" + update.InvoiceId},
		},
		UpdateExpression: aws.String(`set #k1 = :v1, #k2 = :v2, #k3 = :v3, #k4 = :v4, #k5 = :v5,
		 #k6 = :v6, #k7 = :v7, #k8 = :v8, #k9 = :v9, #k10 = :v10, #k11 = :v11, #k12 = :v12`),
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
			"#k12": "invoiceNumber",
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
			":v10": &types.AttributeValueMemberL{Value: serviceList},
			":v11": &types.AttributeValueMemberS{Value: update.VendorEmail},
			":v12": &types.AttributeValueMemberS{Value: update.InvoiceNumber},
		},
	})
	return updateError
}

func UpdateInvoiceStatuses(client *dynamodb.Client, tableName string, orgId string, invoiceIds []string, status string) error {
	transactInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{},
	}

	for _, invoiceId := range invoiceIds {
		transactInput.TransactItems = append(transactInput.TransactItems, types.TransactWriteItem{
			Update: &types.Update{
				TableName: &tableName,
				Key: map[string]types.AttributeValue{
					"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
					"sk": &types.AttributeValueMemberS{Value: models.INVOICE + "#" + invoiceId},
				},
				UpdateExpression: aws.String("set #k1 = :v1, #k2 = :v2"),
				ExpressionAttributeNames: map[string]string{
					"#k1": "status",
					"#k2": "lsi1sk",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":v1": &types.AttributeValueMemberS{Value: status},
					":v2": &types.AttributeValueMemberS{Value: models.STATUS + "#" + status},
				},
			},
		})

		if len(transactInput.TransactItems) == TRANSACTION_LIMIT {
			_, transactionError := client.TransactWriteItems(context.TODO(), &transactInput)
			if transactionError != nil {
				utils.LogError("error while writing batches to db", transactionError)
				return transactionError
			}
			transactInput.TransactItems = []types.TransactWriteItem{}
		}
	}

	_, transactionError := client.TransactWriteItems(context.TODO(), &transactInput)
	if transactionError != nil {
		utils.LogError("error while writing batches to db", transactionError)
		return transactionError
	}

	return nil
}

func UpdateVendor(client *dynamodb.Client, tableName string, orgId string, update VendorUpdate) error {
	_, updateError := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			"sk": &types.AttributeValueMemberS{Value: models.VENDOR + "#" + update.VendorName},
		},
		UpdateExpression: aws.String("set #pk = :pk, #sk = :sk, #vendorEmail = :vendorEmail"),
		ExpressionAttributeNames: map[string]string{
			"#pk":          "pk",
			"#sk":          "sk",
			"#vendorEmail": "vendorEmail",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":          &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			":sk":          &types.AttributeValueMemberS{Value: models.VENDOR + "#" + update.VendorName},
			":vendorEmail": &types.AttributeValueMemberS{Value: update.VendorEmail},
		},
	})
	return updateError
}
