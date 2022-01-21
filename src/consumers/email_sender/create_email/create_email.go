package create_email

import (
	"bytes"
	"context"
	"html/template"
	"path"
	"pocok/src/utils/models"
	"runtime"

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

type emailTemplateData struct {
	ApiUrl string
}

func GetHtmlSummary(apiUrl string) (string, error) {
	_, filename, _, _ := runtime.Caller(0)
	currentPath := path.Dir(filename)

	t, err := template.ParseFiles(currentPath + "/../../../amp/email-summary.html")
	if err != nil {
		return "", err
	}
	var templateBuffer bytes.Buffer

	templateData := emailTemplateData{
		ApiUrl: apiUrl,
	}
	execerr := t.Execute(&templateBuffer, templateData)
	if execerr != nil {
		return "", err
	}
	return templateBuffer.String(), nil
}
