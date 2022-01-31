package db

import (
	"context"
	"pocok/src/utils"
	"pocok/src/utils/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func UpdateInvoiceStatus(client *dynamodb.Client, tableName string, orgId string, invoiceId string, status string) error {
	_, err := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
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
	})
	return err
}

func UpdateInvoiceStatuses(client *dynamodb.Client, tableName string, orgId string, invoiceIds []string, status string) error {
	var batchInput dynamodb.BatchWriteItemInput

	for _, invoiceId := range invoiceIds {
		batchInput.RequestItems[tableName] = append(batchInput.RequestItems[tableName], types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					"pk":     &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
					"sk":     &types.AttributeValueMemberS{Value: models.INVOICE + "#" + invoiceId},
					"status": &types.AttributeValueMemberS{Value: status},
					"lsi1sk": &types.AttributeValueMemberS{Value: models.STATUS + "#" + status},
				},
			},
		})

		if len(batchInput.RequestItems[tableName]) == 25 {
			_, err := client.BatchWriteItem(context.TODO(), &batchInput)
			if err != nil {
				utils.LogError("error while writing batches to db", err)
				return err
			}
			batchInput.RequestItems[tableName] = []types.WriteRequest{}
		}
	}

	_, err := client.BatchWriteItem(context.TODO(), &batchInput)
	if err != nil {
		utils.LogError("error while writing batches to db", err)
		return err
	}

	return nil
}
