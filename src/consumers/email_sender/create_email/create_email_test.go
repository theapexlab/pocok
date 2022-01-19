package create_email_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "pocok/src/cron/invoice_summary/create_email"
	"pocok/src/mock"
	"pocok/src/utils/models"
)

var _ = Describe("CreateEmail", func() {
	var invoices []models.Invoice
	var email *models.Email
	var err error

	When("invoices are empty", func() {
		BeforeEach(func() {
			invoices = []models.Invoice{}
			email, err = CreateEmail(invoices)
		})

		It("returns a valid email response", func() {
			Expect(email).ToNot(BeNil())
		})

		It("returns the empty email template", func() {
			Expect(email.Html).ToNot(Equal(""))
		})

		It("has no attachments", func() {
			Expect(len(email.Attachments)).To(Equal(0))
		})

		It("does not error ", func() {
			Expect(err).To(BeNil())
		})
	})

	When("invoices are not empty", func() {
		BeforeEach(func() {
			invoices = mock.Invoices
			email, err = CreateEmail(invoices)
		})

		It("returns a valid email response", func() {
			Expect(email).ToNot(BeNil())
		})

		It("returns a proper email response", func() {
			Expect(email.Html).ToNot(Equal(""))
		})

		It("all invoices have attachments", func() {
			Expect(len(email.Attachments)).To(Equal(len(invoices)))
		})

		It("does not error ", func() {
			Expect(err).To(BeNil())
		})
	})
})
