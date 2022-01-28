package db

import (
	"context"
	"errors"
	"pocok/src/utils"
	"pocok/src/utils/models"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func DeleteInvoice(dbClient *dynamodb.Client, tableName string, s3Client s3.Client, bucketName string, orgId string, invoiceId string) error {
	_, s3Err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &invoiceId,
	})
	if s3Err != nil {
		utils.LogError("Error while removing object from s3 bucket", s3Err)
	}

	_, dbErr := dbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			"sk": &types.AttributeValueMemberS{Value: models.INVOICE + "#" + invoiceId},
		},
	})
	if dbErr != nil {
		utils.LogError("Error while removing object from dynamoDB", dbErr)
	}

	if s3Err != nil || dbErr != nil {
		return errors.New(dbErr.Error() + s3Err.Error())
	}
	return nil
}
