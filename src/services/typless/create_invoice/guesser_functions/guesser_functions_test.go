package guesser_functions_test

import (
	"pocok/src/mocks/typless/parse_mock_json"
	"pocok/src/services/typless"
	"pocok/src/services/typless/create_invoice/guesser_functions"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Guesser functions", func() {
	var extractedData *typless.ExtractDataFromFileOutput

	When("2021-000022.json is received as an input", func() {
		BeforeEach(func() {
			extractedData = parse_mock_json.Parse("2021-000022.json")
		})

		It("guesses vendor name succesfully", func() {
			vendorName := guesser_functions.GuessVendorName(&extractedData.TextBlocks)
			Expect(vendorName).To(Equal("Teszt Elek (Kisadózó)"))
		})
		It("guesses invoice number succesfully", func() {
			originalFilename := "2021-000022.pdf"
			invoiceNumber := guesser_functions.GuessInvoiceNumberFromFilename(originalFilename, &extractedData.TextBlocks)
			Expect(invoiceNumber).To(Equal("2021-000022"))
		})

		It("guesses iban succesfully", func() {
			iban := guesser_functions.GuessIban(&extractedData.TextBlocks)
			Expect(iban).To(Equal("HU93116000060000000012345676"))
		})

		It("guesses hungarian bank account number succesfully", func() {
			bankAccountNumber := guesser_functions.GuessHunBankAccountNumber(&extractedData.TextBlocks)
			Expect(bankAccountNumber).To(Equal("12345678-10589326-49010011"))
		})

		It("guesses currency type succesfully", func() {
			currency := guesser_functions.GuessCurrency(&extractedData.TextBlocks)
			Expect(currency).To(Equal("HUF"))
		})

		It("guesses due date succesfully", func() {
			dueDate := guesser_functions.GuessDueDate(&extractedData.TextBlocks)
			Expect(dueDate).To(Equal("2021-01-02"))
		})
	})

	When("E0041.json is received as an input", func() {
		BeforeEach(func() {
			extractedData = parse_mock_json.Parse("E0041.json")
		})

		It("guesses vendor name succesfully", func() {
			vendorName := guesser_functions.GuessVendorName(&extractedData.TextBlocks)
			Expect(vendorName).To(Equal("Endrödi Zoltán EV"))
		})

		It("guesses invoice number succesfully", func() {
			originalFilename := "endrodi_E0041.pdf"
			invoiceNumber := guesser_functions.GuessInvoiceNumberFromFilename(originalFilename, &extractedData.TextBlocks)
			Expect(invoiceNumber).To(Equal("E0041"))
		})

		It("guesses iban succesfully", func() {
			iban := guesser_functions.GuessIban(&extractedData.TextBlocks)
			Expect(iban).To(Equal(""))
		})

		It("guesses hungarian bank account number succesfully", func() {
			bankAccountNumber := guesser_functions.GuessHunBankAccountNumber(&extractedData.TextBlocks)
			Expect(bankAccountNumber).To(Equal("11736006-20410614"))
		})

		It("guesses currency type succesfully", func() {
			currency := guesser_functions.GuessCurrency(&extractedData.TextBlocks)
			Expect(currency).To(Equal("HUF"))
		})

		It("guesses due date succesfully", func() {
			dueDate := guesser_functions.GuessDueDate(&extractedData.TextBlocks)
			Expect(dueDate).To(Equal("2021-11-30"))
		})

	})
	When("SZERV-2021-42.json is received as an input", func() {
		BeforeEach(func() {
			extractedData = parse_mock_json.Parse("SZERV-2021-42.json")
		})

		It("guesses vendor name succesfully", func() {
			vendorName := guesser_functions.GuessVendorName(&extractedData.TextBlocks)
			Expect(vendorName).To(Equal("SZERVEZET- ÉS"))
		})

		It("guesses invoice number succesfully", func() {
			originalFilename := "SZERV-2021-42.pdf"
			invoiceNumber := guesser_functions.GuessInvoiceNumberFromFilename(originalFilename, &extractedData.TextBlocks)
			Expect(invoiceNumber).To(Equal("SZERV-2021-42"))
		})

		It("guesses iban succesfully", func() {
			iban := guesser_functions.GuessIban(&extractedData.TextBlocks)
			Expect(iban).To(Equal(""))
		})

		It("guesses hungarian bank account number succesfully", func() {
			bankAccountNumber := guesser_functions.GuessHunBankAccountNumber(&extractedData.TextBlocks)
			Expect(bankAccountNumber).To(Equal(""))
		})

		It("guesses currency type succesfully", func() {
			currency := guesser_functions.GuessCurrency(&extractedData.TextBlocks)
			Expect(currency).To(Equal("HUF"))
		})

		It("guesses due date succesfully", func() {
			dueDate := guesser_functions.GuessDueDate(&extractedData.TextBlocks)
			Expect(dueDate).To(Equal("2021-09-23"))
		})

	})

	When("oszp.json is received as an input", func() {
		BeforeEach(func() {
			extractedData = parse_mock_json.Parse("oszp.json")
		})

		It("guesses vendor name succesfully", func() {
			vendorName := guesser_functions.GuessVendorName(&extractedData.TextBlocks)
			Expect(vendorName).To(Equal("TEST BETÉTI TÁRSASÁG"))
		})

		It("guesses invoice number succesfully", func() {
			originalFilename := "TEST-2021-42.pdf"
			invoiceNumber := guesser_functions.GuessInvoiceNumberFromFilename(originalFilename, &extractedData.TextBlocks)
			Expect(invoiceNumber).To(Equal("TEST-2021-42"))
		})

		It("guesses iban succesfully", func() {
			iban := guesser_functions.GuessIban(&extractedData.TextBlocks)
			Expect(iban).To(Equal(""))
		})

		It("guesses hungarian bank account number succesfully", func() {
			bankAccountNumber := guesser_functions.GuessHunBankAccountNumber(&extractedData.TextBlocks)
			Expect(bankAccountNumber).To(Equal("12345678-12345678-12345678"))
		})

		It("guesses currency type succesfully", func() {
			currency := guesser_functions.GuessCurrency(&extractedData.TextBlocks)
			Expect(currency).To(Equal("HUF"))
		})

		It("guesses due date succesfully", func() {
			dueDate := guesser_functions.GuessDueDate(&extractedData.TextBlocks)
			Expect(dueDate).To(Equal("2021-09-23"))
		})

	})
	When("billingo.json is received as an input", func() {
		BeforeEach(func() {
			extractedData = parse_mock_json.Parse("billingo.json")
		})

		It("guesses vendor name incorrectly", func() {
			vendorName := guesser_functions.GuessVendorName(&extractedData.TextBlocks)
			Expect(vendorName).To(Equal("ELEKTRONIKUS SZAMLA"))
		})

		It("guesses invoice number succesfully", func() {
			originalFilename := "A-2021-111.pdf"
			invoiceNumber := guesser_functions.GuessInvoiceNumberFromFilename(originalFilename, &extractedData.TextBlocks)
			Expect(invoiceNumber).To(Equal("A-2021-111"))
		})

		It("guesses iban succesfully", func() {
			iban := guesser_functions.GuessIban(&extractedData.TextBlocks)
			Expect(iban).To(Equal("HU19120105010040405600200005"))
		})

		It("guesses hungarian bank account number succesfully", func() {
			bankAccountNumber := guesser_functions.GuessHunBankAccountNumber(&extractedData.TextBlocks)
			Expect(bankAccountNumber).To(Equal("12345678-12345678-12345678"))
		})

		It("guesses currency type succesfully", func() {
			currency := guesser_functions.GuessCurrency(&extractedData.TextBlocks)
			Expect(currency).To(Equal("HUF"))
		})

		It("guesses due date succesfully", func() {
			dueDate := guesser_functions.GuessDueDate(&extractedData.TextBlocks)
			Expect(dueDate).To(Equal("2021-12-07"))
		})

	})
})
