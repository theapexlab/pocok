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
		s3Resp, s3Error := client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    &invoice.Filename,
		})
		if s3Error != nil {
			return attachments, s3Error
		}
		file, readError := ioutil.ReadAll(s3Resp.Body)
		if readError != nil {
			return attachments, readError
		}
		attachments[invoice.Filename] = file
	}
	return attachments, nil
}

type emailTemplateData struct {
	ApiUrl    string
	Token     string
	Accepted  string
	Rejected  string
	PocokLogo string
}

func GetHtmlSummary(apiUrl string, logoUrl string) (string, error) {
	summaryTemplate, summaryError := summary_email_template.Get()

	if summaryError != nil {
		utils.LogError("Error while reading in summary file", summaryError)
		return "", summaryError
	}

	t, templateError := template.New("Template").Delims("[[", "]]").Parse(summaryTemplate)
	if templateError != nil {
		utils.LogError("Error while creating template.", templateError)
		return "", templateError
	}

	token, tokenError := auth.CreateToken(models.APEX_ID)
	if tokenError != nil {
		utils.LogError("Error while creating token.", tokenError)
		return "", tokenError
	}

	templateData := emailTemplateData{
		ApiUrl:    apiUrl,
		Token:     token,
		Accepted:  models.ACCEPTED,
		Rejected:  models.REJECTED,
		PocokLogo: logoUrl,
	}

	var templateBuffer bytes.Buffer
	executionError := t.Execute(&templateBuffer, templateData)

	if executionError != nil {
		utils.LogError("Error while executing template insetion.", executionError)
		return "", executionError
	}
	return templateBuffer.String(), nil
}
