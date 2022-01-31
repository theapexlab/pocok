package mailgun

import (
	"github.com/mailgun/mailgun-go/v4"
)

func GetClient(domain string, apiKey string) *mailgun.MailgunImpl {
	client := mailgun.NewMailgun(domain, apiKey)
	client.SetAPIBase(mailgun.APIBaseEU)

	return client
}
