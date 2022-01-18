package create_email

import (
	"bytes"
	"html/template"
	"pocok/src/utils/models"
)

func GetHtmlSummary(invoices []models.Invoice) (string, error) {
	ids := make([]string, len(invoices))
	for i, inv := range invoices {
		ids[i] = inv.Id
	}

	t, err := template.ParseFiles("../amp/email-summary.html")
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

func GetHtmlSummaryEmpty() (string, error) {
	t, err := template.ParseFiles("../amp/email-summary-empty.html")
	if err != nil {
		return "", err
	}
	var templateBuffer bytes.Buffer
	execerr := t.Execute(&templateBuffer, nil)
	if execerr != nil {
		return "", err
	}
	return templateBuffer.String(), nil
}

func GetAttachments(invoices []models.Invoice) []string {
	attachments := []string{}
	for _, invoice := range invoices {
		attachments = append(attachments, invoice.Filename)
	}
	return attachments
}

func CreateEmail(invoices []models.Invoice) (*models.Email, error) {
	to := "billing@apexlab.io"
	subject := "Pocok Invoice Summary"

	var html string
	var err error
	if len(invoices) == 0 {
		html, err = GetHtmlSummaryEmpty()
	} else {
		html, err = GetHtmlSummary(invoices)
	}

	if err != nil {
		return nil, err
	}
	attachments := GetAttachments(invoices)

	email := models.Email{
		To:          to,
		Subject:     subject,
		Html:        html,
		Attachments: attachments,
	}
	return &email, nil
}
