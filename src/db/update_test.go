package db_test

import (
	"pocok/src/db"
	"pocok/src/utils/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Status update validation", func() {
	var update db.StatusUpdate
	var testError error

	When("the input is valid", func() {
		BeforeEach(func() {
			update, testError = db.CreateStatusUpdate(map[string]string{
				"invoiceId": "ID1",
				"status":    models.ACCEPTED,
			})
		})

		It("does not error", func() {
			Expect(testError).To(BeNil())
		})

		It("contains the data", func() {
			Expect(update.InvoiceId).To(Equal("ID1"))
			Expect(update.Status).To(Equal(models.ACCEPTED))
		})
	})

	When("the invoiceId is missing", func() {
		BeforeEach(func() {
			_, testError = db.CreateStatusUpdate(map[string]string{
				"status": models.ACCEPTED,
			})
		})

		It("errors", func() {
			Expect(testError).ToNot(BeNil())
		})
	})

	When("the status is missing", func() {
		BeforeEach(func() {
			_, testError = db.CreateStatusUpdate(map[string]string{
				"invoiceId": "asd",
			})
		})

		It("errors", func() {
			Expect(testError).ToNot(BeNil())
		})
	})

	When("the the status is invalid", func() {
		BeforeEach(func() {
			_, testError = db.CreateStatusUpdate(map[string]string{
				"status":    models.PENDING,
				"invoiceId": "ID1",
			})
		})

		It("errors", func() {
			Expect(testError).ToNot(BeNil())
		})
	})

})
