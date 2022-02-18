package parse_email_test

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"pocok/src/api/process_email/parse_email"
	"pocok/src/mocks/raw_emails/read_raw_email"
	"pocok/src/utils/models"
)

var _ = Describe("ParseEmail", func() {
	contentUrl := "https://pipedream-emails.s3.amazonaws.com/88vrfi98nk1mf6qnmknp26qkbiuvsvc37s59h9g1?AWSAccessKeyId=AKIA5F5AGIEASBWKVUEZ&Expires=1645179119&Signature=5FZiJZaRJKQ0nTd%2F5FDnTwUbBJw%3D"
	var rawEmail string
	var invoiceMessages []models.UploadInvoiceMessage
	var testError error

	BeforeEach(func() {
		httpmock.Activate()
	})

	AfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	When("body is malformed", func() {
		BeforeEach(func() {
			invoiceMessages, testError = parse_email.ParseEmail("")
		})

		It("returns nil", func() {
			Expect(invoiceMessages).To(BeNil())
		})

		It("errors", func() {
			Expect(testError).To(MatchError(models.ErrInvalidJson))
		})
	})

	When("body doesn't contain pdf attachment", func() {
		BeforeEach(func() {
			rawEmail = read_raw_email.Read("without_attachement.txt")

			httpmock.Reset()
			httpmock.RegisterResponder("GET", contentUrl, httpmock.NewStringResponder(200, rawEmail))

			invoiceMessages, testError = parse_email.ParseEmail(`{
				"mail": {
					"content_url": "` + contentUrl + `"
				}
			}`)
		})

		It("returns nil", func() {
			Expect(invoiceMessages).To(BeNil())
		})

		It("errors", func() {
			Expect(testError).To(MatchError(parse_email.ErrNoPdfAttachmentFound))
		})
	})

	When("body does contain a pdf attachment", func() {
		BeforeEach(func() {
			rawEmail = read_raw_email.Read("with_attachement.txt")

			httpmock.Reset()
			httpmock.RegisterResponder("GET", contentUrl, httpmock.NewStringResponder(200, rawEmail))

			invoiceMessages, testError = parse_email.ParseEmail(`{
				"mail": {
					"content_url": "` + contentUrl + `"
				}
			}`)
		})

		It("not errors", func() {
			Expect(testError).To(BeNil())
		})

		It("returns invoice", func() {
			invoiceMessage := invoiceMessages[0]
			Expect(invoiceMessage.Type).To(Equal("base64"))
			Expect(invoiceMessage.Filename).To(Equal("test-pdf-invoice.pdf"))
		})
	})

	When("body contains multiple pdf attachments", func() {
		BeforeEach(func() {
			rawEmail = read_raw_email.Read("multiple_attachements.txt")

			httpmock.Reset()
			httpmock.RegisterResponder("GET", contentUrl, httpmock.NewStringResponder(200, rawEmail))

			invoiceMessages, testError = parse_email.ParseEmail(`{
				"mail": {
					"content_url": "` + contentUrl + `"
				}
			}`)
		})

		It("not errors", func() {
			Expect(testError).To(BeNil())
		})

		It("returns the invoices correctly", func() {
			Expect(len(invoiceMessages)).To(Equal(2))
			Expect(invoiceMessages[0].Type).To(Equal("base64"))
			Expect(invoiceMessages[0].Filename).To(Equal("test-pdf-invoice-copy.pdf"))

			Expect(invoiceMessages[1].Type).To(Equal("base64"))
			Expect(invoiceMessages[1].Filename).To(Equal("test-pdf-invoice.pdf"))
		})
	})
})
