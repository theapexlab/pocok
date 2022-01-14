package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"pocok/src/utils/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/cavaliergopher/grab/v3"
	"github.com/segmentio/ksuid"
)

type dependencies struct {
	bucketName      string
	snsTopicArn     string
	textractRoleArn string
	s3Client        *s3.Client
	textractClient  *textract.Client
}

func main() {
	d := &dependencies{
		bucketName:      os.Getenv("bucketName"),
		snsTopicArn:     os.Getenv("snsTopicArn"),
		textractRoleArn: os.Getenv("textractRoleArn"),
		s3Client:        aws_clients.GetS3Client(),
		textractClient:  aws_clients.GetTextractClient(),
	}

	lambda.Start(d.handler)
}

func (d *dependencies) handler(event events.SQSEvent) error {
	for _, record := range event.Records {
		uploadInvoiceMessage, err := parseBody(record.Body)
		if err != nil {
			continue
		}

		uploadPDFErr := uploadPDF(d, uploadInvoiceMessage)
		// if the original file doesn't exists, no need to retry the message
		if uploadPDFErr != nil && uploadPDFErr != grab.StatusCodeError(403) {
			return uploadPDFErr
		}
	}

	return nil
}

func parseBody(body string) (*models.UploadInvoiceMessage, error) {
	var jsonBody *models.UploadInvoiceMessage

	if err := json.Unmarshal([]byte(body), &jsonBody); err != nil {
		return nil, err
	}

	return jsonBody, nil
}

func uploadPDF(d *dependencies, uploadInvoiceMessage *models.UploadInvoiceMessage) error {
	var data []byte
	var err error

	switch uploadInvoiceMessage.Type {
	case "url":
		data, err = downloadFile(uploadInvoiceMessage.Body)
	case "base64":
		data, err = base64.StdEncoding.DecodeString(uploadInvoiceMessage.Body)
	default:
		err = errors.New("invalid uploadInvoiceMessage type: " + uploadInvoiceMessage.Type)
	}

	if err != nil {
		utils.LogError("", err)
		return err
	}

	filename := ksuid.New().String() + ".pdf"

	_, s3Err := d.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &d.bucketName,
		Key:         &filename,
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/pdf"),
	})

	if s3Err != nil {
		utils.LogError("Error while uploading to s3", s3Err)
		return s3Err
	}

	_, textractErr := d.textractClient.StartExpenseAnalysis(context.TODO(), &textract.StartExpenseAnalysisInput{
		DocumentLocation: &types.DocumentLocation{
			S3Object: &types.S3Object{
				Bucket: &d.bucketName,
				Name:   &filename,
			},
		},
		NotificationChannel: &types.NotificationChannel{
			SNSTopicArn: &d.snsTopicArn,
			RoleArn:     &d.textractRoleArn,
		},
	})
	if textractErr != nil {
		utils.LogError("Error while starting textract", textractErr)

		_, s3Err := d.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: &d.bucketName,
			Key:    &filename,
		})
		if s3Err != nil {
			utils.LogError("Error while deleting file from s3", s3Err)
		}

		return textractErr
	}

	return nil
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
