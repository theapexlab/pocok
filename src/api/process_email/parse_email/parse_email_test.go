package parse_email_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "pocok/src/api/process_email/parse_email"
	"pocok/src/utils/models"
)

var _ = Describe("ParseEmail", func() {
	var invoice *models.UploadInvoiceMessage
	var err error

	When("body is malformed", func() {
		BeforeEach(func() {
			invoice, err = ParseEmail("")
		})

		It("returns nil", func() {
			Expect(invoice).To(BeNil())
		})

		It("errors", func() {
			Expect(err).To(MatchError(models.ErrInvalidJson))
		})
	})

	When("body doesn't contain pdf attachment", func() {
		BeforeEach(func() {
			invoice, err = ParseEmail(`{
				"attachments": [
					{
						"contentType": "image/gif",
						"content_b64": "lkjasdlfkjasdlfkjasldfkjasldkfjasdfkl",
						"length": "37",
						"transferEncoding": "base64"
					}
				]
			}`)
		})

		It("returns nil", func() {
			Expect(invoice).To(BeNil())
		})

		It("errors", func() {
			Expect(err).To(MatchError(ErrNoPdfAttachmentFound))
		})
	})

	When("body does contain a pdf attachment", func() {
		BeforeEach(func() {
			invoice, err = ParseEmail(`{
				"attachments": [
					{
						"contentType": "application/pdf",
						"content_b64": "lkjasdlfkjasdlfkjasldfkjasldkfjasdfkl",
						"length": "37",
						"transferEncoding": "base64"
					}
				]
			}`)
		})

		It("not errors", func() {
			Expect(err).To(BeNil())
		})

		It("returns invoice", func() {
			Expect(invoice.Type).To(Equal("base64"))
			Expect(invoice.Body).To(Equal("lkjasdlfkjasdlfkjasldfkjasldkfjasdfkl"))
		})
	})
})
