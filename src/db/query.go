package db

import (
	"context"
	"pocok/src/utils"
	"pocok/src/utils/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetPendingInvoices(client *dynamodb.Client, tableName string) ([]models.Invoice, error) {
	// TODO query the pending ones
	resp, err := client.BatchGetItem(context.TODO(), &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			// "Invoices": types.KeysAndAttributes{},
		},
	})
	if err != nil {
		utils.LogError("Error while querying the db", err)
		return []models.Invoice{}, err
	}

	invoiceTable := resp.Responses["Invoices"]

	invoices := []models.Invoice{}
	for _, item := range invoiceTable {
		invoice := models.Invoice{}
		err := attributevalue.UnmarshalMap(item, invoice)
		if err != nil {
			utils.LogError("Error while loading invoices", err)
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}
