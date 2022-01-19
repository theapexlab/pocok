package create_email

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"os"
	"pocok/src/utils/models"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetAttachments(client *s3.Client, bucketName string, invoices []models.Invoice) (map[string][]byte, error) {
	attachments := map[string][]byte{}
	for _, invoice := range invoices {
		s3Resp, s3Err := client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    &invoice.Filename,
		})
		if s3Err != nil {
			return attachments, s3Err
		}
		file := []byte{}
		_, err := s3Resp.Body.Read(file)
		if err != nil {
			return attachments, err
		}
		attachments[invoice.Filename] = file
	}
	return attachments, nil
}

func GetHtmlSummary(invoices []models.Invoice) (string, error) {
	ids := make([]string, len(invoices))
	for i, inv := range invoices {
		ids[i] = inv.InvoiceId
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.ParseFiles(wd + "/src/amp/email-summary.html")
	if err != nil {
		return "", err
	}
	var templateBuffer bytes.Buffer
	execerr := t.Execute(&templateBuffer, ids)
	if execerr != nil {
		return "", err
	}
	return templateBuffer.String(), nil
}
