package aws_clients

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/textract"
)

func GetSQSClient() *sqs.Client {
	cfg, getSqsClientError := config.LoadDefaultConfig(context.TODO())
	if getSqsClientError != nil {
		panic("SQS configuration error, " + getSqsClientError.Error())
	}

	return sqs.NewFromConfig(cfg)
}

func GetS3Client() *s3.Client {
	cfg, getS3ClientError := config.LoadDefaultConfig(context.TODO())
	if getS3ClientError != nil {
		panic("S3 configuration error, " + getS3ClientError.Error())
	}

	return s3.NewFromConfig(cfg)
}

func GetDbClient() *dynamodb.Client {
	cfg, getDbClientError := config.LoadDefaultConfig(context.TODO())
	if getDbClientError != nil {
		panic("DynamoDB configuration error, " + getDbClientError.Error())
	}

	return dynamodb.NewFromConfig(cfg)
}

func GetTextractClient() *textract.Client {
	cfg, getTextractClientError := config.LoadDefaultConfig(context.TODO())
	if getTextractClientError != nil {
		panic("Textract configuration error, " + getTextractClientError.Error())
	}

	return textract.NewFromConfig(cfg)
}
