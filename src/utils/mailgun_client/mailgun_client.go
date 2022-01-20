package mailgun_client

import (
	"os"

	"github.com/mailgun/mailgun-go/v4"
)

func GetMailgunClient() *mailgun.MailgunImpl {
	domain := os.Getenv("mailgunDomain")
	apiKey := os.Getenv("mailgunApiKey")
	client := mailgun.NewMailgun(domain, apiKey)
	return client
}
