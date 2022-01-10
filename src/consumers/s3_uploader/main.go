package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"pocok/src/utils/aws_clients"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cavaliergopher/grab/v3"
	"github.com/segmentio/ksuid"
)

type dependencies struct {
	tableName  string
	bucketName string
	s3Client   *s3.Client
	dbClient   *dynamodb.Client
}

func downloadFile(url string) ([]byte, error) {
	client := grab.NewClient()
	req, err := grab.NewRequest(".", url)
	if err != nil {
		return []byte(nil), err
	}

	// store file in memory
	req.NoStore = true
	resp := client.Do(req)
	data, err := resp.Bytes()

	return data, err
}

func uploadPDF(d *dependencies, url string) error {
	data, downloadErr := downloadFile(url)
	if downloadErr != nil {
		fmt.Printf("❌ Error while downloading PDF: %s", downloadErr)
		return downloadErr
	}

	filename := strconv.FormatInt(time.Now().Unix(), 10) + ".pdf"

	s3Resp, s3Err := d.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &d.bucketName,
		Key:    &filename,
		Body:   bytes.NewReader(data),
	})
	if s3Err != nil {
		fmt.Printf("❌ Error while uploading to s3: %s", s3Err)
		return s3Err
	}

	_, dbErr := d.dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item: map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: ksuid.New().String()},
			"filename": &types.AttributeValueMemberS{Value: filename},
			"etag":     &types.AttributeValueMemberS{Value: *s3Resp.ETag},
		},
	})
	if dbErr != nil {
		fmt.Printf("❌ Error while creating record in DB: %s", dbErr)
		return dbErr
	}

	return nil
}

func (d *dependencies) handler(event events.SQSEvent) error {
	for _, record := range event.Records {
		url := record.Body
		uploadPDFErr := uploadPDF(d, url)
		// if the original file doesn't exists, no need to retry the message
		if uploadPDFErr != nil && uploadPDFErr != grab.StatusCodeError(403) {
			return uploadPDFErr
		}
	}

	return nil
}

func main() {
	d := &dependencies{
		tableName:  os.Getenv("tableName"),
		bucketName: os.Getenv("bucketName"),
		s3Client:   aws_clients.GetS3Client(),
		dbClient:   aws_clients.GetDbClient(),
	}

	lambda.Start(d.handler)
}
