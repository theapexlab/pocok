package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/cavaliergopher/grab/v3"
)

type dependencies struct {
	bucketName             string
	processInvoiceQueueUrl string
	sqsClient              *sqs.Client
	s3Client               *s3.Client
	tableName              string
	dbClient               *dynamodb.Client
}

func main() {
	d := &dependencies{
		processInvoiceQueueUrl: os.Getenv("processInvoiceQueueUrl"),
		s3Client:               aws_clients.GetS3Client(),
		sqsClient:              aws_clients.GetSQSClient(),
		tableName:              os.Getenv("tableName"),
		dbClient:               aws_clients.GetDbClient(),
		bucketName:             os.Getenv("bucketName"),
	}

	lambda.Start(d.handler)
}

func (d *dependencies) handler(event events.SQSEvent) error {
	for _, record := range event.Records {
		uploadInvoiceMessage, parseError := parseBody(record.Body)
		if parseError != nil {
			continue
		}

		uploadPDFError := uploadPDF(d, uploadInvoiceMessage)
		// if the original file doesn't exist, no need to retry the message
		if uploadPDFError != nil && uploadPDFError != grab.StatusCodeError(403) {
			return uploadPDFError
		}
	}

	return nil
}

func parseBody(body string) (*models.UploadInvoiceMessage, error) {
	var jsonBody *models.UploadInvoiceMessage

	if unmarshalError := json.Unmarshal([]byte(body), &jsonBody); unmarshalError != nil {
		return nil, unmarshalError
	}

	return jsonBody, nil
}

func uploadPDF(d *dependencies, uploadInvoiceMessage *models.UploadInvoiceMessage) error {
	var data []byte
	var uploadError error

	switch uploadInvoiceMessage.Type {
	case "url":
		data, uploadError = downloadFile(uploadInvoiceMessage.Body)
	case "base64":
		data, uploadError = base64.StdEncoding.DecodeString(uploadInvoiceMessage.Body)
	default:
		uploadError = errors.New("invalid uploadInvoiceMessage type: " + uploadInvoiceMessage.Type)
	}

	if uploadError != nil {
		utils.LogError("", uploadError)
		return uploadError
	}

	checksum := sha256.Sum256(data)
	checksumString := fmt.Sprintf("%x", checksum)
	filename := checksumString + ".pdf"

	s3Response, s3LoadError := d.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &d.bucketName,
		Key:    &filename,
	})

	// if the file already exists, no need to continue
	if s3Response != nil {
		utils.Log("invoice already exists!")
		return nil
	}

	var nck *types.NoSuchKey

	if s3LoadError != nil && !errors.As(s3LoadError, &nck) {
		utils.LogError("s3 network error!", s3LoadError)
		return s3LoadError
	}

	_, s3UploadError := d.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &d.bucketName,
		Key:         &filename,
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/pdf"),
		Metadata: map[string]string{
			"OriginalFilename": uploadInvoiceMessage.Filename,
		},
	})

	if s3UploadError != nil {
		utils.LogError("Error while uploading to s3", s3UploadError)
		return s3UploadError
	}

	_, sqsError := d.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(filename),
		QueueUrl:    &d.processInvoiceQueueUrl,
	})

	if sqsError != nil {
		utils.LogError("Error while sending message to ProcessInvoice queue", sqsError)

		_, s3Error := d.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: &d.bucketName,
			Key:    &filename,
		})
		if s3Error != nil {
			utils.LogError("Error while deleting file from s3", s3Error)
		}

		return sqsError
	}

	return nil
}

func downloadFile(url string) ([]byte, error) {
	client := grab.NewClient()
	req, requestError := grab.NewRequest(".", url)
	if requestError != nil {
		return []byte(nil), requestError
	}

	// store file in memory
	req.NoStore = true
	resp := client.Do(req)
	data, responseError := resp.Bytes()

	return data, responseError
}
