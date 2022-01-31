package db

import (
	"context"
	"errors"
	"pocok/src/utils"
	"pocok/src/utils/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type StatusUpdate struct {
	InvoiceId string
	Status    string
}
type GenericUpdate struct {
	InvoiceId string
}

func CreateStatusUpdate(data map[string]string) (StatusUpdate, error) {
	var patch StatusUpdate
	err := utils.MapToStruct(data, &patch)
	if err != nil {
		return patch, errors.New("invalid input")
	}
	if patch.InvoiceId == "" {
		return patch, errors.New("invalid invoiceId")
	}
	if patch.Status != models.ACCEPTED && patch.Status != models.REJECTED {
		return patch, errors.New("invalid status")
	}
	return patch, nil
}

func UpdateInvoiceStatus(client *dynamodb.Client, tableName string, orgId string, update StatusUpdate) error {
	_, err := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			"sk": &types.AttributeValueMemberS{Value: models.INVOICE + "#" + update.InvoiceId},
		},
		UpdateExpression: aws.String("set #k1 = :v1, #k2 = :v2"),
		ExpressionAttributeNames: map[string]string{
			"#k1": "status",
			"#k2": "lsi1sk",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: update.Status},
			":v2": &types.AttributeValueMemberS{Value: models.STATUS + "#" + update.Status},
		},
	})
	return err
}

func UpdateInvoice(client *dynamodb.Client, tableName string, update GenericUpdate) error {
	return nil
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
