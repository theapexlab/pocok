package parse_email

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"pocok/src/api/process_email/url_parsing_strategies"
	"pocok/src/utils"
	"pocok/src/utils/models"

	"github.com/DusanKasan/parsemail"
)

var ErrNoPdfAttachmentFound = errors.New("no pdf attachment found")
var ErrEmailFromSenderAddress = errors.New("email sent from mailgun sender address")

func ParseEmail(body string) ([]models.UploadInvoiceMessage, error) {
	var jsonBody models.EmailWebhookBody

	if unmarshalError := json.Unmarshal([]byte(body), &jsonBody); unmarshalError != nil {
		return nil, models.ErrInvalidJson
	}

	if jsonBody.Mail.ContentUrl != "" {
		rawEmail, downloadErr := utils.DownloadFile(jsonBody.Mail.ContentUrl)
		if downloadErr != nil {
			utils.LogError("error while downloading raw email", downloadErr)
			return nil, downloadErr
		}

		reader := bytes.NewReader(rawEmail)
		email, parseEmailErr := parsemail.Parse(reader)
		if parseEmailErr != nil {
			utils.LogError("error while parsing raw email", parseEmailErr)
			return nil, parseEmailErr
		}

		if isSentFromRecipientAddress(&email) {
			return nil, ErrEmailFromSenderAddress
		}

		messages := []models.UploadInvoiceMessage{}

		for _, attachment := range email.Attachments {
			if attachment.ContentType == "application/pdf" {
				buffer := new(bytes.Buffer)
				_, readFromErr := buffer.ReadFrom(attachment.Data)
				if readFromErr != nil {
					utils.LogError("error while reading attachment data", readFromErr)
					return nil, readFromErr
				}
				base64Data := base64.StdEncoding.EncodeToString(buffer.Bytes())

				message := models.UploadInvoiceMessage{
					Type:     "base64",
					Body:     base64Data,
					Filename: attachment.Filename,
				}

				messages = append(messages, message)
			}
		}

		if ok, url := hasPdfUrl(&email); ok {
			messages := []models.UploadInvoiceMessage{{
				Type: "url",
				Body: url,
			}}
			return messages, nil
		}

		if len(messages) == 0 {
			return nil, ErrNoPdfAttachmentFound
		}

		return messages, nil
	}

	return nil, ErrNoPdfAttachmentFound
}

func isSentFromRecipientAddress(email *parsemail.Email) bool {
	return len(email.From) > 0 && email.From[0].Address == os.Getenv("mailgunSender")
}

func hasPdfUrl(email *parsemail.Email) (bool, string) {
	url, parseError := url_parsing_strategies.GetPdfUrlFromEmail(email)
	if !errors.Is(parseError, url_parsing_strategies.ErrNoUrlParsingStrategyFound) {
		utils.LogError("error while parsing url from email", parseError)
	}
	return parseError == nil && url != "", url
}
