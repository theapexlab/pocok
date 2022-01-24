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

func GetPendingInvoices(client *dynamodb.Client, tableName string) ([]models.Invoice, error) {
	resp, err := client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              &tableName,
		IndexName:              aws.String(models.LOCAL_SECONDARY_INDEX_1),
		KeyConditionExpression: aws.String("#PK = :PK and #SK = :SK"),
		ExpressionAttributeNames: map[string]string{
			"#PK": "pk",
			"#SK": "lsi1sk",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":PK": &types.AttributeValueMemberS{Value: "ORG#" + models.APEX_ID},
			":SK": &types.AttributeValueMemberS{Value: "STATUS#pending"},
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
