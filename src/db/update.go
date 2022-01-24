package db

import (
	"context"
	"pocok/src/utils/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func UpdateInvoiceStatus(client *dynamodb.Client, tableName string, orgId string, invoiceId string, status string) error {
	_, err := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			"SK": &types.AttributeValueMemberS{Value: models.INVOICE + "#" + invoiceId},
		},
		UpdateExpression: aws.String("set #k1 = :v1 and #k2 = :v2"),
		ExpressionAttributeNames: map[string]string{
			"k1": "status",
			"k2": "lsi1sk",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: status},
			":v2": &types.AttributeValueMemberS{Value: models.STATUS + "#" + status},
		},
	})
	return err
}
