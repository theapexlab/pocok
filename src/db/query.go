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
	resp, err := client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              &tableName,
		IndexName:              aws.String(models.LOCAL_SECONDARY_INDEX_1),
		KeyConditionExpression: aws.String("#PK = :PK and #SK = :SK"),
		ExpressionAttributeNames: map[string]string{
			"#PK": "pk",
			"#SK": "lsi1sk",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":PK": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			":SK": &types.AttributeValueMemberS{Value: models.STATUS + "#pending"},
		},
	})
	if err != nil {
		utils.LogError("Error while querying the db", err)
		return []models.Invoice{}, err
	}

	invoices := []models.Invoice{}
	for _, item := range resp.Items {
		invoice := models.Invoice{}
		err := attributevalue.UnmarshalMap(item, &invoice)
		if err != nil {
			utils.LogError("Error while loading invoices", err)
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}

func GetInvoice(client *dynamodb.Client, tableName string, orgId string, invoiceId string) (*models.Invoice, error) {
	resp, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			"sk": &types.AttributeValueMemberS{Value: models.INVOICE + "#" + invoiceId},
		},
	})
	if err != nil {
		utils.LogError("Error while getting invoice from db", err)
		return nil, err
	}

	invoice := models.Invoice{}
	err = attributevalue.UnmarshalMap(resp.Item, &invoice)
	if err != nil {
		utils.LogError("Error while loading invoice", err)
		return nil, err
	}

	return &invoice, nil
}

func GetVendor(client *dynamodb.Client, tableName string, orgId string, vendorName string) (*models.Vendor, error) {
	resp, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			"sk": &types.AttributeValueMemberS{Value: models.VENDOR + "#" + vendorName},
		},
	})
	if err != nil {
		utils.LogError("Error while getting vendor from db", err)
		return nil, err
	}

	vendor := models.Vendor{}
	err = attributevalue.UnmarshalMap(resp.Item, &vendor)
	if err != nil {
		utils.LogError("Error while loading vendor", err)
		return nil, err
	}

	return &vendor, nil
}
