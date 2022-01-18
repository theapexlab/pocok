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
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("SQS configuration error, " + err.Error())
	}

	return sqs.NewFromConfig(cfg)
}

func GetS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("S3 configuration error, " + err.Error())
	}

	return s3.NewFromConfig(cfg)
}

func GetDbClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("DynamoDB configuration error, " + err.Error())
	}

	return dynamodb.NewFromConfig(cfg)
}

func GetTextractClient() *textract.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Textract configuration error, " + err.Error())
	}

	return textract.NewFromConfig(cfg)
}
