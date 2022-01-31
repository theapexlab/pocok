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

type VendorUpdate struct {
	VendorName  string
	VendorEmail string
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
