package create_invoice_test

import (
	"pocok/src/mocks/typless/parse_mock_json"
	"pocok/src/services/typless"
	"pocok/src/services/typless/create_invoice"
	"pocok/src/utils/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateInvoice", func() {
	var createInvoiceService *create_invoice.CreateInvoiceService
	var extractedData *typless.ExtractDataFromFileOutput
	var invoice *models.Invoice
	var testError error

	When("recieves billingo invoice with normal fields", func() {
		BeforeEach(func() {
			extractedData = parse_mock_json.Parse("billingo.json")
			createInvoiceService = &create_invoice.CreateInvoiceService{
				OriginalFilename: "A-1984-145.pdf",
				ExtractedData:    extractedData,
			}
			invoice = createInvoiceService.CreateInvoice()
		})

		It("not errors", func() {
			Expect(testError).To(BeNil())
		})

		It("return invoice with correct fields", func() {
			Expect(invoice).NotTo(BeNil())

			Expect(invoice.InvoiceNumber).To(Equal("A-2021-111"))
			Expect(invoice.Iban).To(Equal("HU19-123456789101112130200005"))
			Expect(invoice.AccountNumber).To(Equal("HU40 12345678-12345678-12345678"))
			Expect(invoice.VendorName).To(Equal("dr.Test Zoltán ev."))
			Expect(invoice.DueDate).To(Equal("2021-12-07"))
			Expect(invoice.GrossPrice).To(Equal("580.5000"))
			Expect(invoice.Currency).To(Equal("€"))
			Expect(invoice.VatRate).To(Equal("AAM"))
			Expect(invoice.NetPrice).To(Equal("580.5000"))

			Expect(len(invoice.Services)).To(Equal(1))
			Expect(invoice.Services[0].Name).To(Equal("Óradíjas tanácsadás"))
			Expect(invoice.Services[0].Unit).To(Equal("óra"))
			Expect(invoice.Services[0].Amount).To(Equal("13.5000"))
			Expect(invoice.Services[0].NetPrice).To(Equal("580.5000"))
			Expect(invoice.Services[0].GrossPrice).To(Equal("580.5000"))
			Expect(invoice.Services[0].UnitNetPrice).To(Equal(""))
			Expect(invoice.Services[0].VatRate).To(Equal(""))
			Expect(invoice.Services[0].VatAmount).To(Equal(""))

		})
	})
	When("recieves Online számlázo program invoice with multiple line items", func() {
		BeforeEach(func() {
			extractedData = parse_mock_json.Parse("oszp.json")
			createInvoiceService = &create_invoice.CreateInvoiceService{
				OriginalFilename: "TEST-2021-42.pdf",
				ExtractedData:    extractedData,
			}
			invoice = createInvoiceService.CreateInvoice()
		})

		It("not errors", func() {
			Expect(testError).To(BeNil())
		})

		It("return invoice with correct main fields", func() {
			Expect(invoice).NotTo(BeNil())

			Expect(invoice.InvoiceNumber).To(Equal("TEST-2021-42"))
			Expect(invoice.Iban).To(Equal(""))
			Expect(invoice.AccountNumber).To(Equal("12345678-12345678-12345678"))
			Expect(invoice.VendorName).To(Equal("TEST BETÉTI TÁRSASÁG"))
			Expect(invoice.DueDate).To(Equal("2021-09-23"))
			Expect(invoice.GrossPrice).To(Equal("24800.0000"))
			Expect(invoice.Currency).To(Equal("HUF"))
			Expect(invoice.VatRate).To(Equal("AAM"))
			Expect(invoice.NetPrice).To(Equal("24800.0000"))
		})

		It("return invoice with correct line item fields", func() {
			Expect(len(invoice.Services)).To(Equal(2))
			Expect(invoice.Services[0].Name).To(Equal("standup"))
			Expect(invoice.Services[0].Unit).To(Equal("Ora / Hour"))
			Expect(invoice.Services[0].Amount).To(Equal("2.0000"))
			Expect(invoice.Services[0].NetPrice).To(Equal("23600.0000"))
			Expect(invoice.Services[0].GrossPrice).To(Equal("23600.0000"))
			Expect(invoice.Services[0].UnitNetPrice).To(Equal(""))
			Expect(invoice.Services[0].VatRate).To(Equal(""))
			Expect(invoice.Services[0].VatAmount).To(Equal(""))

			Expect(invoice.Services[1].Name).To(Equal("travel costs"))
			Expect(invoice.Services[1].Unit).To(Equal("Darab / Piece"))
			Expect(invoice.Services[1].Amount).To(Equal("1.0000"))
			Expect(invoice.Services[1].NetPrice).To(Equal("4000.0000"))
			Expect(invoice.Services[1].GrossPrice).To(Equal("4000.0000"))
			Expect(invoice.Services[1].UnitNetPrice).To(Equal(""))
			Expect(invoice.Services[1].VatRate).To(Equal(""))
			Expect(invoice.Services[1].VatAmount).To(Equal(""))

		})
	})
	When("recieves Online számlázo program invoice with no service description line item", func() {
		BeforeEach(func() {
			extractedData = parse_mock_json.Parse("oszp.json")
			extractedData.LineItems[0][3].Values[0].Value = ""
			createInvoiceService = &create_invoice.CreateInvoiceService{
				OriginalFilename: "TEST-2021-42.pdf",
				ExtractedData:    extractedData,
			}
			invoice = createInvoiceService.CreateInvoice()
		})

		It("not errors", func() {
			Expect(testError).To(BeNil())
		})

		It("skips line item with no service description", func() {
			Expect(invoice).NotTo(BeNil())
			Expect(len(invoice.Services)).To(Equal(1))

		})
	})
	When("recieves Online számlázo program invoice with no grossprice, netprice and currency line item", func() {
		BeforeEach(func() {
			extractedData = parse_mock_json.Parse("oszp.json")
			extractedData.LineItems[0][1].Values[0].Value = ""
			extractedData.LineItems[0][2].Values[0].Value = ""
			createInvoiceService = &create_invoice.CreateInvoiceService{
				OriginalFilename: "TEST-2021-42.pdf",
				ExtractedData:    extractedData,
			}
			invoice = createInvoiceService.CreateInvoice()
		})

		It("not errors", func() {
			Expect(testError).To(BeNil())
		})

		It("skips line item with empty field no grossprice, netprice and currency", func() {
			Expect(invoice).NotTo(BeNil())
			Expect(len(invoice.Services)).To(Equal(1))

		})
	})
})
