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

func DeleteInvoice(dbClient *dynamodb.Client, tableName string, s3Client s3.Client, bucketName string, orgId string, invoiceId string, filename string) error {
	_, s3Error := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &filename,
	})
	if s3Error != nil {
		utils.LogError("Error while removing object from s3 bucket", s3Error)
	}

	_, dbError := dbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + orgId},
			"sk": &types.AttributeValueMemberS{Value: models.INVOICE + "#" + invoiceId},
		},
	})
	if dbError != nil {
		utils.LogError("Error while removing object from dynamoDB", dbError)
	}

	if s3Error != nil || dbError != nil {
		return errors.New(dbError.Error() + s3Error.Error())
	}
	return nil
}
