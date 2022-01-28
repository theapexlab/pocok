package create_training_data_test

import (
	"pocok/src/mocks"
	"pocok/src/services/typless"
	"pocok/src/services/typless/create_training_data"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func getLearningFieldValue(learningFields []typless.LearningField, field string) string {
	for _, f := range learningFields {
		if f.Name == field {
			return f.Value
		}
	}

	return ""
}

var _ = Describe("CreateTrainingData", func() {
	var trainingData *typless.TrainingData

	When("gets invoice with all fields", func() {
		BeforeEach(func() {
			trainingData = create_training_data.CreateTrainingData(&mocks.MockInvoice)
		})

		It("returns fields with correct values", func() {
			Expect(trainingData.DocumentObjectId).To(Equal("0e809bfab6a4253a1e1cfdfa5088d30380565c02"))

			Expect(getLearningFieldValue(trainingData.LearningFields, typless.INVOICE_NUMBER)).To(Equal("500000"))
			Expect(getLearningFieldValue(trainingData.LearningFields, typless.GROSS_PRICE)).To(Equal("20000"))
			Expect(getLearningFieldValue(trainingData.LearningFields, typless.NET_PRICE)).To(Equal("10000"))
			Expect(getLearningFieldValue(trainingData.LearningFields, typless.VENDOR_NAME)).To(Equal("Csipkés Zoltán"))
			Expect(getLearningFieldValue(trainingData.LearningFields, typless.ACCOUNT_NUMBER)).To(Equal("10001000-10001000-10001000"))
			Expect(getLearningFieldValue(trainingData.LearningFields, typless.IBAN)).To(Equal("HU69119800810030005009212644"))
			Expect(getLearningFieldValue(trainingData.LearningFields, typless.DUE_DATE)).To(Equal("2050.01.01."))
			Expect(getLearningFieldValue(trainingData.LearningFields, typless.CURRENCY)).To(Equal("huf"))
			Expect(getLearningFieldValue(trainingData.LearningFields, typless.VAT_RATE)).To(Equal("27%"))
			Expect(getLearningFieldValue(trainingData.LearningFields, typless.VAT_AMOUNT)).To(Equal("2700"))

			Expect(getLearningFieldValue(trainingData.LineItems[0], typless.SERVICE_NAME)).To(Equal("Kutya"))
			Expect(getLearningFieldValue(trainingData.LineItems[0], typless.SERVICE_UNIT)).To(Equal("db"))
			Expect(getLearningFieldValue(trainingData.LineItems[0], typless.SERVICE_AMOUNT)).To(Equal("500"))
			Expect(getLearningFieldValue(trainingData.LineItems[0], typless.SERVICE_NET_PRICE)).To(Equal("5000"))
			Expect(getLearningFieldValue(trainingData.LineItems[0], typless.SERVICE_GROSS_PRICE)).To(Equal("10000"))
			Expect(getLearningFieldValue(trainingData.LineItems[0], typless.SERVICE_CURRENCY)).To(Equal("huf"))
			Expect(getLearningFieldValue(trainingData.LineItems[0], typless.SERVICE_VAT_RATE)).To(Equal("27%"))
			Expect(getLearningFieldValue(trainingData.LineItems[0], typless.SERVICE_VAT_AMOUNT)).To(Equal("2700"))
		})
	})

	When("gets invoice with missing fields", func() {
		BeforeEach(func() {
			trainingData = create_training_data.CreateTrainingData(&mocks.MockInvoiceMissingFields)
		})

		It("returns fields with correct values", func() {
			Expect(getLearningFieldValue(trainingData.LearningFields, typless.INVOICE_NUMBER)).To(Equal("2012-12"))
			Expect(getLearningFieldValue(trainingData.LearningFields, typless.GROSS_PRICE)).To(Equal(""))

			Expect(len(trainingData.LineItems)).To(Equal(0))
		})
	})
})
