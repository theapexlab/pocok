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
	Filename  string
}
type GenericUpdate struct {
	InvoiceId string
}

const TRANSACTION_LIMIT = 25 // 25 is the max number of items that can be updated in a single transaction

type VendorUpdate struct {
	VendorName  string
	VendorEmail string
}

func CreateStatusUpdate(data map[string]string) (*StatusUpdate, error) {
	var update StatusUpdate
	err := utils.MapToStruct(data, &update)
	if err != nil {
		return nil, errors.New("invalid input")
	}
	if update.InvoiceId == "" {
		return nil, errors.New("invalid invoiceId")
	}
	if update.Status != models.ACCEPTED && update.Status != models.REJECTED {
		return nil, errors.New("invalid status")
	}
	if update.Status == models.REJECTED && update.Filename == "" {
		return nil, errors.New("invalid update, filename must be present on reject")
	}
	return &update, nil
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
	transactInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{},
	}

	for _, invoiceId := range invoiceIds {
		transactInput.TransactItems = append(transactInput.TransactItems, types.TransactWriteItem{
			Update: &types.Update{
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
			},
		})

		if len(transactInput.TransactItems) == TRANSACTION_LIMIT {
			_, err := client.TransactWriteItems(context.TODO(), &transactInput)
			if err != nil {
				utils.LogError("error while writing batches to db", err)
				return err
			}
			transactInput.TransactItems = []types.TransactWriteItem{}
		}
	}

	_, err := client.TransactWriteItems(context.TODO(), &transactInput)
	if err != nil {
		utils.LogError("error while writing batches to db", err)
		return err
	}

	return nil
}

func UpdateVendor(client *dynamodb.Client, tableName string, orgId string, update VendorUpdate) error {
	_, err := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			"sk": &types.AttributeValueMemberS{Value: models.VENDOR + "#" + update.VendorName},
		},
		UpdateExpression: aws.String("set #pk = :pk, #sk = :sk, #vendorEmail = :vendorEmail"),
		ExpressionAttributeNames: map[string]string{
			"#pk":          "pk",
			"#sk":          "sk",
			"#vendorEmail": "vendorEmail",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":          &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			":sk":          &types.AttributeValueMemberS{Value: models.VENDOR + "#" + update.VendorName},
			":vendorEmail": &types.AttributeValueMemberS{Value: update.VendorEmail},
		},
	})
	return err
}
