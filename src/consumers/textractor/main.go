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
		documentTextDetectionMessage, err := parseBody(record.Body)
		if err != nil {
			utils.LogError("Failed to parse textract message", err)
			return err
		}

		// see job status types: &types.JobStatus.Values()
		if documentTextDetectionMessage.Status != "SUCCEEDED" {
			err := errors.New("text detection status failed")
			utils.LogError("Status != SUCCEEDED", err)
			return err
		}

		fmt.Println(documentTextDetectionMessage) //todo: remove this line

		// todo: call GetDocumentTextDetection api with JobId
		res, detectErr := d.textractClient.GetDocumentTextDetection(context.TODO(), &textract.GetDocumentTextDetectionInput{
			JobId: &documentTextDetectionMessage.JobId,
		})

		if detectErr != nil {
			utils.LogError("GetDocumentTextDetection failed", detectErr)
			return detectErr

			// AccessDeniedException: User: arn:aws:sts::382372657671:assumed-role/arn-aws-iam-382372657671-textractQueueConsumertex-1474UXEUAJT5Z/arn-aws-iam-382372657671--textractQueueConsumertex-iNBB0FV6hwHB
			// is not authorized to perform: textract:GetDocumentTextDetection because no identity-based policy allows the textract:GetDocumentTextDetection action
		}

		fmt.Println(res) //todo: remove this line
	}

	return nil
}

func parseBody(body string) (*documentTextDetectionMessage, error) {
	var jsonBody *documentTextDetectionBody
	var jsonMessage *documentTextDetectionMessage

	if err := json.Unmarshal([]byte(body), &jsonBody); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(jsonBody.Message), &jsonMessage); err != nil {
		return nil, err
	}

	return jsonMessage, nil
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
