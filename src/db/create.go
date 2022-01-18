package db

import (
	"context"
	"pocok/src/utils/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func PutInvoice(client *dynamodb.Client, tableName string, invoice models.Invoice) (*dynamodb.PutItemOutput, error) {
	item, itemErr := attributevalue.MarshalMap(invoice)
	if itemErr != nil {
		return nil, itemErr
	}

	dbResp, dbErr := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})
	return dbResp, dbErr
}
