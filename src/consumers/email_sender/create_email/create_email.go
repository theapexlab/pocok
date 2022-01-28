package create_email

import (
	"bytes"
	"context"
	"html/template"
	"io/ioutil"
	"pocok/src/amp/summary_email_template"
	"pocok/src/utils"
	"pocok/src/utils/auth"
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
		file, readErr := ioutil.ReadAll(s3Resp.Body)
		if readErr != nil {
			return attachments, readErr
		}
		attachments[invoice.Filename] = file
	}
	return attachments, nil
}

type emailTemplateData struct {
	ApiUrl   string
	Token    string
	Accepted string
	Rejected string
}

func GetHtmlSummary(apiUrl string) (string, error) {
	t, templateErr := template.New("Template").Delims("[[", "]]").Parse(summary_email_template.Get())
	if templateErr != nil {
		utils.LogError("Error while creating template.", templateErr)
		return "", templateErr
	}

	token, tokenErr := auth.CreateToken(models.APEX_ID)
	if tokenErr != nil {
		utils.LogError("Error while creating token.", tokenErr)
		return "", tokenErr
	}
	templateData := emailTemplateData{
		ApiUrl:   apiUrl,
		Token:    token,
		Accepted: models.ACCEPTED,
		Rejected: models.REJECTED,
	}

	var templateBuffer bytes.Buffer
	executionErr := t.Execute(&templateBuffer, templateData)

	if executionErr != nil {
		utils.LogError("Error while executing template insetion.", executionErr)
		return "", executionErr
	}
	return templateBuffer.String(), nil
}
