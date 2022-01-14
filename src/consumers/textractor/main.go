package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/textract"
)

// TODO: WIP implementation of POCO-25: create pdf data extract consumer

type dependencies struct {
	bucketName     string
	s3Client       *s3.Client
	textractClient *textract.Client
	dbClient       *dynamodb.Client
}

type documentTextDetectionBody struct {
	Message string `json:"Message"`
}
type documentTextDetectionMessage struct {
	JobId  string `json:"JobId,omitempty"`
	Status string `json:"Status,omitempty"`
}

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
		documentTextDetectionMessage, err := parseRecordBody(record.Body)
		if err != nil {
			utils.LogError("Failed to parse textract message", err)
			return err
		}

		// FYI: see job status types: &types.JobStatus.Values()
		if documentTextDetectionMessage.Status != "SUCCEEDED" {
			err := errors.New("text detection status failed")
			utils.LogError("Status != SUCCEEDED", err)
			return err
		}

		fmt.Println(documentTextDetectionMessage) //todo: remove this line

		// todo: call GetDocumentTextDetection api with JobId
		d.getResults(documentTextDetectionMessage)

		// todo:  save attrubutes to db after textract completed
		// _, dbErr := d.dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		// 	TableName: &d.tableName,
		// 	Item: map[string]types.AttributeValue{
		// 		"id":       &types.AttributeValueMemberS{Value: ksuid.New().String()},
		// 		"filename": &types.AttributeValueMemberS{Value: filename},
		// 	},
		// })

		// if dbErr != nil {
		// 	utils.LogError("Error while inserting to db", dbErr)
		// 	return dbErr
		// }
	}

	return nil
}

func parseRecordBody(event string) (*documentTextDetectionMessage, error) {
	var snsEntity events.SNSEntity

	if err := json.Unmarshal([]byte(event), &snsEntity); err != nil {
		return nil, err
	}

	var message *documentTextDetectionMessage

	if err := json.Unmarshal([]byte(snsEntity.Message), &message); err != nil {
		return nil, err
	}

	return message, nil
}

func (d *dependencies) getResults(message *documentTextDetectionMessage) {
	res, err := d.textractClient.GetDocumentTextDetection(context.TODO(), &textract.GetDocumentTextDetectionInput{
		JobId: &message.JobId,
	})
	if err != nil {
		utils.LogError("", err)
	}

	fmt.Println(res)

	// for _, block := range res.Blocks {
	// 	block.
	// }
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

// func getS3Object(d *dependencies, filename string) (*s3.GetObjectOutput, error) {
// 	return d.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
// 		Bucket: &d.bucketName,
// 		Key:    &filename,
// 	})
// }
