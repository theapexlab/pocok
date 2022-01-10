package main

import (
	"bytes"
	"context"
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
	tableName string
	s3Bucket  string
	s3Client  *s3.Client
	dbClient  *dynamodb.Client
}

func downloadFile(url string) ([]byte, error) {
	client := grab.NewClient()
	req, err := grab.NewRequest(".", url)
	if err != nil {
		return []byte(nil), err
	}
	req.NoStore = true
	resp := client.Do(req)
	data, err := resp.Bytes()
	return data, err
}

func uploadPDF(d *dependencies, url string) error {
	data, _ := downloadFile(url)
	filename := strconv.FormatInt(time.Now().Unix(), 10) + ".pdf"

	s3Resp, err := d.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &d.s3Bucket,
		Key:    &filename,
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}

	_, err = d.dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item: map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: ksuid.New().String()},
			"filename": &types.AttributeValueMemberS{Value: filename},
			"etag":     &types.AttributeValueMemberS{Value: *s3Resp.ETag},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *dependencies) handler(event events.SQSEvent) {
	for _, record := range event.Records {
		url := record.Body
		uploadPDF(d, url)
	}
}

func main() {
	d := &dependencies{
		tableName: os.Getenv("tableName"),
		s3Bucket:  os.Getenv("bucketName"),
		s3Client:  aws_clients.GetS3Client(),
		dbClient:  aws_clients.GetDbClient(),
	}

	lambda.Start(d.handler)
}
