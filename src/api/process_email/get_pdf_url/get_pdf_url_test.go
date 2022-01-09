package get_pdf_url_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "pocok/src/api/process_email/get_pdf_url"
)

var _ = Describe("GetPdfUrl", func() {
	var pdfUrl string
	var err error

	When("body is malformed", func() {
		BeforeEach(func() {
			pdfUrl, err = GetPdfUrl("")
		})

		It("returns empty string", func() {
			Expect(pdfUrl).To(Equal(""))
		})

		It("errors", func() {
			Expect(err).To(MatchError(ErrInvalidJson))
		})
	})

	When("body doesn't contain pdf attachment", func() {
		BeforeEach(func() {
			pdfUrl, err = GetPdfUrl(`{"attachment-1":{"filename":"test.gif","encoding":"7bit","mimetype":"image/gif","url":"https://s3.amazonaws.com/test-bucket/test.gif"}}`)
		})

		It("returns empty string", func() {
			Expect(pdfUrl).To(Equal(""))
		})

		It("errors", func() {
			Expect(err).To(MatchError(ErrNoPdfAttachmentFound))
		})
	})

	When("body is valid", func() {
		BeforeEach(func() {
			pdfUrl, err = GetPdfUrl(`{"attachment-1":{"filename":"test-pdf-invoice.pdf","encoding":"7bit","mimetype":"application/pdf","url":"https://s3.amazonaws.com/test-bucket/test-pdf-invoice.pdf"}}`)
		})

		It("returns pdf url", func() {
			Expect(pdfUrl).To(Equal("https://s3.amazonaws.com/test-bucket/test-pdf-invoice.pdf"))
		})

		It("does not error", func() {
			Expect(err).To(BeNil())
		})
	})
})
