package db

import (
	"context"
	"errors"
	"pocok/src/utils"
	"pocok/src/utils/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/segmentio/ksuid"
)

type PutInvoiceInput struct {
	OrgId   string
	Invoice models.Invoice
}

func ValidatePutInvoiceInput(input PutInvoiceInput) error {
	if input.OrgId == "" {
		return errors.New("orgId empty")
	}
	if input.Invoice.Filename == "" {
		return errors.New("invoice filename empty")
	}
	return nil
}

func PutInvoice(client *dynamodb.Client, tableName string, input PutInvoiceInput) (*dynamodb.PutItemOutput, error) {
	validationErr := ValidatePutInvoiceInput(input)
	if validationErr != nil {
		utils.LogError("input validation error", validationErr)
		return nil, validationErr
	}

	id := ksuid.New().String()
	status := models.PENDING
	reveicedAt := time.Now().Format(time.RFC3339)

	invoice := &models.Invoice{
		Pk:              models.ORG + "#" + input.OrgId,
		Sk:              models.INVOICE + "#" + id,
		Lsi1sk:          models.STATUS + "#" + status,
		ReceivedAt:      reveicedAt,
		InvoiceId:       id,
		EntityType:      models.INVOICE,
		Status:          status,
		Filename:        input.Invoice.Filename,
		VendorName:      input.Invoice.VendorName,
		AccountNumber:   input.Invoice.AccountNumber,
		Iban:            input.Invoice.Iban,
		NetPrice:        input.Invoice.NetPrice,
		GrossPrice:      input.Invoice.GrossPrice,
		VatAmount:       input.Invoice.VatAmount,
		VatRate:         input.Invoice.VatRate,
		Currency:        input.Invoice.Currency,
		DueDate:         input.Invoice.DueDate,
		Services:        input.Invoice.Services,
		TyplessObjectId: input.Invoice.TyplessObjectId,
		InvoiceNumber:   input.Invoice.InvoiceNumber,
	}

	vendor, vendorError := GetVendor(client, tableName, GetVendorInput{OrgId: input.OrgId, VendorName: input.Invoice.VendorName})
	if vendorError != nil {
		utils.Log("no vendor with the same name found")
	}
	if vendor != nil {
		invoice.VendorEmail = vendor.VendorEmail
	}
	item, itemError := attributevalue.MarshalMap(invoice)
	if itemError != nil {
		utils.LogError("error while marshaling invoice", itemError)
		return nil, itemError
	}
	dbResp, dbError := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})
	if dbError != nil {
		utils.LogError("error creating item", dbError)
		return nil, dbError
	}

	return dbResp, nil
}
