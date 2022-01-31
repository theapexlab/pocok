package parse_email_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "pocok/src/api/process_email/parse_email"
	"pocok/src/utils/models"
)

var _ = Describe("ParseEmail", func() {
	var invoiceMessage *models.UploadInvoiceMessage
	var err error

	When("body is malformed", func() {
		BeforeEach(func() {
			invoiceMessage, err = ParseEmail("")
		})

		It("returns nil", func() {
			Expect(invoiceMessage).To(BeNil())
		})

		It("errors", func() {
			Expect(err).To(MatchError(models.ErrInvalidJson))
		})
	})

	When("body doesn't contain pdf attachment", func() {
		BeforeEach(func() {
			invoiceMessage, err = ParseEmail(`{
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
			Expect(invoiceMessage).To(BeNil())
		})

		It("errors", func() {
			Expect(err).To(MatchError(ErrNoPdfAttachmentFound))
		})
	})

	When("body does contain a pdf attachment", func() {
		BeforeEach(func() {
			invoiceMessage, err = ParseEmail(`{
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
			Expect(err).To(BeNil())
		})

		It("returns invoice", func() {
			Expect(invoiceMessage.Type).To(Equal("base64"))
			Expect(invoiceMessage.Body).To(Equal("lkjasdlfkjasdlfkjasldfkjasldkfjasdfkl"))
			Expect(invoiceMessage.Filename).To(Equal("SZERV-2021-87.pdf"))
		})
	})
})
