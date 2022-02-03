package parse_email_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"pocok/src/api/process_email/parse_email"
	"pocok/src/utils/models"
)

var _ = Describe("ParseEmail", func() {
	var invoiceMessages []models.UploadInvoiceMessage
	var testError error

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
			invoiceMessages, testError = parse_email.ParseEmail(`{
				"attachments": [
					{
						"contentType": "image/gif",
						"content_b64": "lkjasdlfkjasdlfkjasldfkjasldkfjasdfkl",
						"length": "37",
						"transferEncoding": "base64",
						"fileName": "SZERV-2021-87.pdf"
					}
				]
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
			invoiceMessages, testError = parse_email.ParseEmail(`{
				"attachments": [
					{
						"contentType": "application/pdf",
						"content_b64": "lkjasdlfkjasdlfkjasldfkjasldkfjasdfkl",
						"length": "37",
						"transferEncoding": "base64",
						"fileName": "SZERV-2021-87.pdf"
					}
				]
			}`)
		})

		It("not errors", func() {
			Expect(testError).To(BeNil())
		})

		It("returns invoice", func() {
			invoiceMessage := invoiceMessages[0]
			Expect(invoiceMessage.Type).To(Equal("base64"))
			Expect(invoiceMessage.Body).To(Equal("lkjasdlfkjasdlfkjasldfkjasldkfjasdfkl"))
			Expect(invoiceMessage.Filename).To(Equal("SZERV-2021-87.pdf"))
		})
	})

	When("body contains multiple pdf attachments", func() {
		BeforeEach(func() {
			invoiceMessages, testError = parse_email.ParseEmail(`{
				"attachments": [
					{
						"contentType": "application/pdf",
						"content_b64": "1",
						"length": "1",
						"transferEncoding": "base64",
						"fileName": "BRUH-1.pdf"
					},
					{
						"contentType": "application/pdf",
						"content_b64": "2",
						"length": "1",
						"transferEncoding": "base64",
						"fileName": "BRUH-2.pdf"
					}
				]
			}`)
		})

		It("not errors", func() {
			Expect(testError).To(BeNil())
		})

		It("returns the invoices correctly", func() {
			Expect(len(invoiceMessages)).To(Equal(2))
			Expect(invoiceMessages[0].Type).To(Equal("base64"))
			Expect(invoiceMessages[0].Body).To(Equal("1"))
			Expect(invoiceMessages[0].Filename).To(Equal("BRUH-1.pdf"))

			Expect(invoiceMessages[1].Type).To(Equal("base64"))
			Expect(invoiceMessages[1].Body).To(Equal("2"))
			Expect(invoiceMessages[1].Filename).To(Equal("BRUH-2.pdf"))
		})
	})
})
