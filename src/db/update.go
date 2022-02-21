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

type UpdateStatusInput struct {
	OrgId     string
	InvoiceId string
	Status    string
}
type UpdateDataInput struct {
	OrgId   string
	Invoice models.Invoice
}
type UpdateVendorInput struct {
	OrgId       string
	VendorName  string
	VendorEmail string
}

func ValidateUpdateStatusInput(input UpdateStatusInput) error {
	if input.OrgId == "" {
		return errors.New("orgId empty")
	}
	if input.InvoiceId == "" {
		return errors.New("invoiceId empty")
	}
	if !utils.SliceContains(
		[]string{
			models.ACCEPTED,
			models.TRANSFER_ERROR,
			models.TRANSFER_LOADING,
		},
		input.Status,
	) {
		return errors.New("invalid status")
	}
	return nil
}

func ValidateUpdateDataInput(input UpdateDataInput) error {
	if input.OrgId == "" {
		return errors.New("orgId empty")
	}

	invoice := input.Invoice
	if invoice.InvoiceId == "" {
		return errors.New("invalid invoiceId")
	}

	if invoice.VendorName == "" {
		return errors.New("invalid vendor name")
	}

	if invoice.AccountNumber == "" && invoice.Iban == "" {
		return errors.New("iban or account number must be provided")
	}

	if invoice.Iban != "" {
		_, ibanError := utils.GetValidIban(invoice.Iban)
		if ibanError != nil {
			return ibanError
		}
	}

	if invoice.AccountNumber != "" {
		validAccountNumber, accountNumberError := utils.GetValidAccountNumber(invoice.AccountNumber)
		if accountNumberError != nil {
			return accountNumberError
		}
		invoice.AccountNumber = validAccountNumber
	}

	_, priceError := utils.GetValidPrice(invoice.GrossPrice)
	if priceError != nil {
		return errors.New("invalid gross price")
	}

	_, dateError := utils.GetValidDueDate(invoice.DueDate)
	if dateError != nil {
		return errors.New("invalid due date")
	}

	_, currError := utils.GetValidCurrency(invoice.Currency)
	if currError != nil {
		return currError
	}
	return nil
}

func UpdateInvoiceStatus(client *dynamodb.Client, tableName string, input UpdateStatusInput) error {
	validationErr := ValidateUpdateStatusInput(input)
	if validationErr != nil {
		utils.LogError("input validation error", validationErr)
		return validationErr
	}
	_, updateError := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + input.OrgId},
			"sk": &types.AttributeValueMemberS{Value: models.INVOICE + "#" + input.InvoiceId},
		},
		UpdateExpression: aws.String("set #k1 = :v1, #k2 = :v2"),
		ExpressionAttributeNames: map[string]string{
			"#k1": "status",
			"#k2": "lsi1sk",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: input.Status},
			":v2": &types.AttributeValueMemberS{Value: models.STATUS + "#" + input.Status},
		},
	})
	return updateError
}

func UpdateInvoiceData(client *dynamodb.Client, tableName string, input UpdateDataInput) error {
	validationErr := ValidateUpdateDataInput(input)
	if validationErr != nil {
		utils.LogError("input validation error", validationErr)
		return validationErr
	}
	serviceList, marshalError := attributevalue.MarshalList(input.Invoice.Services)
	if marshalError != nil {
		utils.LogError("Cannot marshal invoice services", marshalError)
		return marshalError
	}

	_, updateError := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + input.OrgId},
			"sk": &types.AttributeValueMemberS{Value: models.INVOICE + "#" + input.Invoice.InvoiceId},
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
			":v1":  &types.AttributeValueMemberS{Value: input.Invoice.VendorName},
			":v2":  &types.AttributeValueMemberS{Value: input.Invoice.AccountNumber},
			":v3":  &types.AttributeValueMemberS{Value: input.Invoice.Iban},
			":v4":  &types.AttributeValueMemberS{Value: input.Invoice.NetPrice},
			":v5":  &types.AttributeValueMemberS{Value: input.Invoice.GrossPrice},
			":v6":  &types.AttributeValueMemberS{Value: input.Invoice.VatAmount},
			":v7":  &types.AttributeValueMemberS{Value: input.Invoice.VatRate},
			":v8":  &types.AttributeValueMemberS{Value: input.Invoice.Currency},
			":v9":  &types.AttributeValueMemberS{Value: input.Invoice.DueDate},
			":v10": &types.AttributeValueMemberL{Value: serviceList},
			":v11": &types.AttributeValueMemberS{Value: input.Invoice.VendorEmail},
			":v12": &types.AttributeValueMemberS{Value: input.Invoice.InvoiceNumber},
		},
	})
	if input.Invoice.VendorName != "" && input.Invoice.VendorEmail != "" {
		updateVendorError := UpdateVendor(client, tableName, UpdateVendorInput{
			OrgId:       input.OrgId,
			VendorName:  input.Invoice.VendorName,
			VendorEmail: input.Invoice.VendorEmail,
		})
		if updateVendorError != nil {
			utils.LogError("error while updating vendor", updateVendorError)
		}
	}

	return updateError
}

func UpdateVendor(client *dynamodb.Client, tableName string, update UpdateVendorInput) error {
	_, updateError := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + update.OrgId},
			"sk": &types.AttributeValueMemberS{Value: models.VENDOR + "#" + update.VendorName},
		},
		UpdateExpression: aws.String("set #vendorEmail = :vendorEmail"),
		ExpressionAttributeNames: map[string]string{
			"#vendorEmail": "vendorEmail",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":vendorEmail": &types.AttributeValueMemberS{Value: update.VendorEmail},
		},
	})
	return updateError
}
