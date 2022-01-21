package create_invoice_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"pocok/src/consumers/invoice_processor/create_invoice"
	"pocok/src/services/typless"
	"pocok/src/utils/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func parseMockJson(filename string) *typless.ExtractDataFromFileOutput {
	var extractedData *typless.ExtractDataFromFileOutput

	mock, readFileErr := ioutil.ReadFile("../../../mocks/typless/" + filename)
	if readFileErr != nil {
		panic("Failed to read mock file")
	}

	if err := json.Unmarshal(mock, &extractedData); err != nil {
		panic("Failed to unmarshal mock file")
	}

	return extractedData
}

var _ = Describe("CreateInvoice", func() {
	var extractedData *typless.ExtractDataFromFileOutput
	var invoice *models.Invoice
	var err error

	When("handles billingo invoice correctly", func() {
		BeforeEach(func() {
			extractedData = parseMockJson("billingo.json")
			invoice = create_invoice.CreateInvoice(extractedData)
		})

		It("not errors", func() {
			Expect(err).To(BeNil())
		})

		It("return invoice with correct fields", func() {
			Expect(invoice).NotTo(BeNil())

			Expect(invoice.InvoiceNumber).To(Equal("E-2021-36"))
			Expect(invoice.Iban).To(Equal("HU19-120105010040405600200005"))
			Expect(invoice.AccountNumber).To(Equal("HU40 12010501-00404056-00100008"))
			Expect(invoice.VendorName).To(Equal("John Doe"))
			Expect(invoice.DueDate).To(Equal("2021. 10.08."))
			Expect(invoice.GrossPrice).To(Equal("322,50"))
			Expect(invoice.NetPrice).To(Equal("322,50"))
			Expect(len(invoice.Services)).To(Equal(1))
			Expect(invoice.Services[0].Name).To(Equal("Óradíjas tanácsadás"))
			Expect(invoice.Services[0].Amount).To(Equal("7.5 óra"))
			Expect(invoice.Services[0].Currency).To(Equal("EUR"))
			Expect(invoice.Services[0].GrossPrice).To(Equal("322,50"))
			Expect(invoice.Services[0].Tax).To(Equal("AAM ÁFA:"))
			fmt.Println(invoice)
		})
	})
})
