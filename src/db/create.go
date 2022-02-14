package db

import (
	"context"
	"pocok/src/utils"
	"pocok/src/utils/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/segmentio/ksuid"
)

func PutInvoice(client *dynamodb.Client, tableName string, invoiceData *models.Invoice) (*dynamodb.PutItemOutput, error) {

	id := ksuid.New().String()
	status := models.PENDING
	reveicedAt := time.Now().Format(time.RFC3339)

	invoice := &models.Invoice{
		Pk:              models.ORG + "#" + models.APEX_ID,
		Sk:              models.INVOICE + "#" + id,
		Lsi1sk:          models.STATUS + "#" + status,
		ReceivedAt:      reveicedAt,
		InvoiceId:       id,
		EntityType:      models.INVOICE,
		Status:          status,
		Filename:        invoiceData.Filename,
		VendorName:      invoiceData.VendorName,
		AccountNumber:   invoiceData.AccountNumber,
		Iban:            invoiceData.Iban,
		NetPrice:        invoiceData.NetPrice,
		GrossPrice:      invoiceData.GrossPrice,
		VatAmount:       invoiceData.VatAmount,
		VatRate:         invoiceData.VatRate,
		Currency:        invoiceData.Currency,
		DueDate:         invoiceData.DueDate,
		Services:        invoiceData.Services,
		TyplessObjectId: invoiceData.TyplessObjectId,
		InvoiceNumber:   invoiceData.InvoiceNumber,
	}

	vendor, vendorError := GetVendor(client, tableName, models.APEX_ID, invoice.VendorName)
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
