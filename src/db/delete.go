package db

import (
	"context"
	"errors"
	"fmt"
	"pocok/src/utils"
	"pocok/src/utils/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type DeleteInvoiceInput struct {
	OrgId     string
	InvoiceId string
}

func GetValidDeleteInvoiceInput(data map[string]string) (*DeleteInvoiceInput, error) {
	var model DeleteInvoiceInput
	mapError := utils.MapToStruct(data, &model)
	if mapError != nil {
		return nil, errors.New("invalid input")
	}
	if model.InvoiceId == "" {
		return nil, errors.New("invoiceId is empty")
	}
	return &model, nil
}

func ValidateDeleteInvoiceInput(input DeleteInvoiceInput) error {
	if input.OrgId == "" {
		return errors.New("orgId empty")
	}
	if input.InvoiceId == "" {
		return errors.New("invoiceId empty")
	}
	return nil
}

func DeleteInvoice(dbClient *dynamodb.Client, tableName string, s3Client s3.Client, bucketName string, input DeleteInvoiceInput) error {
	invoice, invoiceError := GetInvoice(dbClient, tableName, GetInvoiceInput(input))
	fmt.Println(invoice.InvoiceId)
	fmt.Println(invoice.Filename)
	if invoiceError != nil {
		utils.LogError("Error while getting object from db", invoiceError)
		return invoiceError
	}
	_, s3Error := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(invoice.Filename),
	})
	if s3Error != nil {
		utils.LogError("Error while removing object from s3 bucket", s3Error)
		return s3Error
	}

	_, dbError := dbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.ORG + "#" + input.OrgId},
			"sk": &types.AttributeValueMemberS{Value: models.INVOICE + "#" + input.InvoiceId},
		},
	})
	if dbError != nil {
		utils.LogError("Error while removing object from dynamoDB", dbError)
		return dbError
	}

	return nil
}
