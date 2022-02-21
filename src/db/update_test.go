package db_test

import (
	"pocok/src/db"
	"pocok/src/mocks"
	"pocok/src/utils/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Status update validation", func() {
	var updateStatusInput db.UpdateStatusInput
	var updateStatusError error

	When("status is accepted and valid", func() {
		BeforeEach(func() {
			updateStatusInput = db.UpdateStatusInput{
				OrgId:     "ORGID",
				InvoiceId: "INVOICEID",
				Status:    models.ACCEPTED,
			}
			updateStatusError = db.ValidateUpdateStatusInput(updateStatusInput)
		})

		It("does not error", func() {
			Expect(updateStatusError).To(BeNil())
		})
	})

	When("ids are missing", func() {
		BeforeEach(func() {
			updateStatusInput = db.UpdateStatusInput{
				Status: models.PENDING,
			}
			updateStatusError = db.ValidateUpdateStatusInput(updateStatusInput)
		})

		It("errors", func() {
			Expect(updateStatusError).ToNot(BeNil())
		})
	})

	When("status is not valid", func() {
		BeforeEach(func() {
			updateStatusInput = db.UpdateStatusInput{
				OrgId:     "ORGID",
				InvoiceId: "INVOICEID",
				Status:    models.PENDING,
			}
			updateStatusError = db.ValidateUpdateStatusInput(updateStatusInput)
		})

		It("errors", func() {
			Expect(updateStatusError).ToNot(BeNil())
		})
	})
})

var _ = Describe("Data Update valdiation", func() {
	var updateDataInput db.UpdateDataInput
	var updateDataError error

	When("the input is valid", func() {
		BeforeEach(func() {
			invoice := mocks.MockInvoice
			updateDataInput = db.UpdateDataInput{
				OrgId:   "ORGID",
				Invoice: invoice,
			}
			updateDataError = db.ValidateUpdateDataInput(updateDataInput)
		})

		It("errors", func() {
			Expect(updateDataError).To(BeNil())
		})
	})

	When("orgID is empty", func() {
		BeforeEach(func() {
			invoice := mocks.MockInvoice
			updateDataInput = db.UpdateDataInput{
				OrgId:   "",
				Invoice: invoice,
			}
			updateDataError = db.ValidateUpdateDataInput(updateDataInput)
		})

		It("errors", func() {
			Expect(updateDataError).ToNot(BeNil())
		})
	})

	When("invoice id is invalid", func() {
		BeforeEach(func() {
			invoice := mocks.MockInvoice
			invoice.InvoiceId = ""
			updateDataInput = db.UpdateDataInput{
				OrgId:   "ORGID",
				Invoice: invoice,
			}
			updateDataError = db.ValidateUpdateDataInput(updateDataInput)
		})

		It("errors", func() {
			Expect(updateDataError).ToNot(BeNil())
		})
	})

	When("invoice currency is invalid", func() {
		BeforeEach(func() {
			invoice := mocks.MockInvoice
			invoice.Currency = "XDDDD"
			updateDataInput = db.UpdateDataInput{
				OrgId:   "ORGID",
				Invoice: invoice,
			}
			updateDataError = db.ValidateUpdateDataInput(updateDataInput)
		})

		It("errors", func() {
			Expect(updateDataError).ToNot(BeNil())
		})
	})

	When("invoice iban and account numbmer is invalid", func() {
		BeforeEach(func() {
			invoice := mocks.MockInvoice
			invoice.Iban = "XDDDD"
			invoice.AccountNumber = "XDDDD"
			updateDataInput = db.UpdateDataInput{
				OrgId:   "ORGID",
				Invoice: invoice,
			}
			updateDataError = db.ValidateUpdateDataInput(updateDataInput)
		})

		It("errors", func() {
			Expect(updateDataError).ToNot(BeNil())
		})
	})
})
