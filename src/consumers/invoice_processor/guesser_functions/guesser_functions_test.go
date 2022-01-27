package guesser_functions_test

import (
	"pocok/src/consumers/invoice_processor/guesser_functions"
	"pocok/src/mocks/typless/parse_mock_json"
	"pocok/src/services/typless"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Guesser functions", func() {
	var extractedData *typless.ExtractDataFromFileOutput
	// var err error

	When("2021-000022.json is received as an input", func() {
		BeforeEach(func() {
			extractedData = parse_mock_json.Parse("2021-000022.json")
		})

		It("guesses vendor name succesfully", func() {
			vendorName := guesser_functions.GuessVendorName(&extractedData.TextBlocks)
			Expect(vendorName).To(Equal("Teszt Elek"))
		})
		It("guesses invoice number succesfully", func() {
			originalFilename := "2021-000022.pdf"
			invoiceNumber := guesser_functions.GuessInvoiceNumberFromFilename(originalFilename, &extractedData.TextBlocks)
			Expect(invoiceNumber).To(Equal("2021-000022"))
		})

		It("guesses iban succesfully", func() {
			iban := guesser_functions.GuessIbanFromTextBlocks(&extractedData.TextBlocks)
			Expect(iban).To(Equal("HU93116000060000000012345676"))
		})

		It("guesses hungarian bank account number succesfully", func() {
			bankAccountNumber := guesser_functions.GuessHunBankAccountNumberFromTextBlocks(&extractedData.TextBlocks)
			Expect(bankAccountNumber).To(Equal("12345678-10589326-49010011"))
		})

	})
})
