package main

import (
	"fmt"
	"os"
	"pocok/src/utils/aws_clients"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/textract"
)

type dependencies struct {
	bucketName     string
	s3Client       *s3.Client
	textractClient *textract.Client
	dbClient       *dynamodb.Client
}

// func getS3Object(d *dependencies, filename string) (*s3.GetObjectOutput, error) {
// 	return d.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
// 		Bucket: &d.bucketName,
// 		Key:    &filename,
// 	})
// }

func main() {
	d := &dependencies{
		bucketName:     os.Getenv("bucketName"),
		s3Client:       aws_clients.GetS3Client(),
		textractClient: aws_clients.GetTextractClient(),
		dbClient:       aws_clients.GetDbClient(),
	}

	lambda.Start(d.handler)
}

func (d *dependencies) handler(event events.SQSEvent) error {
	for _, record := range event.Records {
		fmt.Println(record.Body)
		// filename := record.Body
		// ParsePdf(d, filename)
	}

	return nil
}

// func ParsePdf(d *dependencies, filename string) error {
// 	fmt.Println("filename", filename)
// 	res, err := d.textractClient.StartDocumentTextDetection (context.TODO(), &textract.AnalyzeExpenseInput{
// 		Document: &types.Document{
// 			S3Object: &types.S3Object{
// 				Bucket: &d.bucketName,
// 				Name:   &filename,
// 			},
// 		},
// 	})

// 	fmt.Println(res)

// 	if err != nil {
// 		utils.LogError("Failed to parse pdf", err)
// 		return err
// 	}

// 	return nil
// }
