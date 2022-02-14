package db

import (
	"context"
	"pocok/src/utils"
	"pocok/src/utils/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetPendingInvoices(client *dynamodb.Client, tableName string, orgId string) ([]models.Invoice, error) {
	resp, dbError := client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              &tableName,
		IndexName:              aws.String(models.LOCAL_SECONDARY_INDEX_1),
		KeyConditionExpression: aws.String("#PK = :PK and begins_with(#SK, :SK)"),
		ExpressionAttributeNames: map[string]string{
			"#PK": "pk",
			"#SK": "lsi1sk",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":PK": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			":SK": &types.AttributeValueMemberS{Value: models.STATUS + "#" + models.PENDING},
		},
	})
	if dbError != nil {
		utils.LogError("Error while querying the db", dbError)
		return []models.Invoice{}, dbError
	}

	invoices := []models.Invoice{}
	for _, item := range resp.Items {
		invoice := models.Invoice{}
		unmarshalError := attributevalue.UnmarshalMap(item, &invoice)
		if unmarshalError != nil {
			utils.LogError("Error while loading invoices", unmarshalError)
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}

func GetInvoice(client *dynamodb.Client, tableName string, orgId string, invoiceId string) (*models.Invoice, error) {
	resp, dbError := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			"sk": &types.AttributeValueMemberS{Value: models.INVOICE + "#" + invoiceId},
		},
	})
	if dbError != nil {
		utils.LogError("Error while getting invoice from db", dbError)
		return nil, dbError
	}

	invoice := models.Invoice{}
	unmarshalError := attributevalue.UnmarshalMap(resp.Item, &invoice)
	if unmarshalError != nil {
		utils.LogError("Error while loading invoice", unmarshalError)
		return nil, unmarshalError
	}

	return &invoice, nil
}

func GetVendor(client *dynamodb.Client, tableName string, orgId string, vendorName string) (*models.Vendor, error) {
	resp, dbError := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			"sk": &types.AttributeValueMemberS{Value: models.VENDOR + "#" + vendorName},
		},
	})
	if dbError != nil {
		utils.LogError("Error while getting vendor from db", dbError)
		return nil, dbError
	}

	vendor := models.Vendor{}
	unmarshalError := attributevalue.UnmarshalMap(resp.Item, &vendor)
	if unmarshalError != nil {
		utils.LogError("Error while loading vendor", unmarshalError)
		return nil, unmarshalError
	}

	return &vendor, nil
}
