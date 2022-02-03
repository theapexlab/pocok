package db_test

import (
	"pocok/src/db"
	"pocok/src/utils/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Status update validation", func() {
	var update *db.StatusUpdate
	var updateError error

	When("status is accept and valid", func() {
		BeforeEach(func() {
			update, updateError = db.CreateStatusUpdate(map[string]string{
				"invoiceId": "ID1",
				"status":    models.ACCEPTED,
			})
		})

		It("does not error", func() {
			Expect(updateError).To(BeNil())
		})

		It("contains the data", func() {
			Expect(update.InvoiceId).To(Equal("ID1"))
			Expect(update.Status).To(Equal(models.ACCEPTED))
		})
	})

	When("id is missing", func() {
		BeforeEach(func() {
			_, updateError = db.CreateStatusUpdate(map[string]string{
				"status": models.ACCEPTED,
			})
		})

		It("errors", func() {
			Expect(updateError).ToNot(BeNil())
		})
	})

	When("status is missing", func() {
		BeforeEach(func() {
			_, updateError = db.CreateStatusUpdate(map[string]string{
				"invoiceId": "asd",
			})
		})

		It("errors", func() {
			Expect(updateError).ToNot(BeNil())
		})
	})

	When("status is pending", func() {
		BeforeEach(func() {
			_, updateError = db.CreateStatusUpdate(map[string]string{
				"status":    models.PENDING,
				"invoiceId": "ID1",
			})
		})

		It("errors", func() {
			Expect(updateError).ToNot(BeNil())
		})
	})

	When("status is reject, but there is no filename", func() {
		BeforeEach(func() {
			_, updateError = db.CreateStatusUpdate(map[string]string{
				"status":    models.REJECTED,
				"invoiceId": "ID1",
			})
		})

		It("errors", func() {
			Expect(updateError).ToNot(BeNil())
		})
	})
})
