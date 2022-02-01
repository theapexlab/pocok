package utils_test

import (
	"pocok/src/mocks"
	"pocok/src/utils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	Describe("MapToStruct", func() {
		type structType struct {
			A string
		}
		var mapData map[string]string
		var structData structType
		var err error

		When("it takes correct params", func() {
			BeforeEach(func() {
				mapData = map[string]string{"a": "cica"}
				structData = structType{}
				err = utils.MapToStruct(mapData, &structData)
			})

			It("does not error", func() {
				Expect(err).To(BeNil())
			})
			It("the struct has the map data", func() {
				Expect(structData.A).To(Equal("cica"))
			})

		})
	})

	Describe("MapInvoiceToInvoiceServiceIndexes", func() {
		When("it receives  an invoice list as param", func() {
			It("returns same invoice values with indexed services", func() {
				invoices := mocks.Invoices
				indexedInvoices := utils.MapInvoiceToInvoiceServiceIndexes(invoices)
				for _, invoice := range indexedInvoices {
					for i, service := range invoice.Services {
						Expect(service.Index).To(Equal(i))
					}
				}
			})
		})
	})

	Describe("MapUpdateDataToInvoice", func() {
		When("it receives valid data map asinput", func() {
			It("returns invoice entity with nested services", func() {
				data := map[string]string{
					"invoiceNumber":          "adfadsf",
					"vendorName":             "adfadsf",
					"vendorEmail":            "adfadsf",
					"accountNumber":          "adfadsf",
					"grossPrice":             "adfadsf",
					"vatAmount":              "adfadsf",
					"vatRate":                "adfadsf",
					"service_name_1":         "Name_1",
					"service_amount_1":       "Amount_1",
					"service_unit_1":         "Unit_1",
					"service_unitNetprice_1": "UnitNetPrice_1",
					"service_netPrice_1":     "NetPrice_1",
					"service_grossPrice_1":   "GrossPrice_1",
					"service_vatRate_1":      "VatAmount_1",
					"service_vatAmount_1":    "VatRate_1",
					"service_currency_1":     "",
				}

				indexedInvoice, err := utils.MapUpdateDataToInvoice(data)

				Expect(err).To(BeNil())

				for i, service := range indexedInvoice.Services {
					Expect(service.Name).To(Equal("Name_" + string(i)))
					Expect(service.Amount).To(Equal("Amount_" + string(i)))
					Expect(service.Unit).To(Equal("Unit_" + string(i)))
					Expect(service.UnitNetPrice).To(Equal("UnitNetPrice_" + string(i)))
					Expect(service.NetPrice).To(Equal("NetPrice_" + string(i)))
					Expect(service.GrossPrice).To(Equal("GrossPrice_" + string(i)))
					Expect(service.VatAmount).To(Equal("VatAmount_" + string(i)))
					Expect(service.VatRate).To(Equal("VatRate_" + string(i)))
				}

			})
		})
	})

	Describe("Validation test", func() {
		When("GetValidAccountNumber receives a valid 16 character long account number", func() {
			It("returns Bank account number and recieves no error", func() {
				mockBankAccNumber := "12345678-12345678-12345678"
				result, err := utils.GetValidAccountNumber(mockBankAccNumber)
				Expect(result).To(Equal("12345678-12345678-12345678"))
				Expect(err).To(BeNil())
			})
		})
		When("GetValidAccountNumber receives a valid 24 character long account number", func() {
			It("returns Bank account number  ", func() {
				mockBankAccNumber := "12345678-12345678"
				result, err := utils.GetValidAccountNumber(mockBankAccNumber)
				Expect(result).To(Equal("12345678-12345678"))
				Expect(err).To(BeNil())
			})
		})

		When("GetValidAccountNumber recieves an invalid account number", func() {
			It("returns an error", func() {
				invalidAcountNumber := "1234-asdf-test"
				_, err := utils.GetValidAccountNumber(invalidAcountNumber)
				Expect(err).To(MatchError("invalid account number"))
			})
		})

		When("GetValidCurrency recieves a currency", func() {
			It("returns converted currency", func() {
				mockCurrency := "HUF"
				result, err := utils.GetValidCurrency(mockCurrency)
				Expect(result).To(Equal("HUF"))
				Expect(err).To(BeNil())
			})
			It("returns converted currency", func() {
				mockCurrency := "Ft"
				result, err := utils.GetValidCurrency(mockCurrency)
				Expect(result).To(Equal("HUF"))
				Expect(err).To(BeNil())
			})
		})

		When("GetValidCurrency recieves a not valid currency", func() {
			It("throws error", func() {
				mockCurrency := "test"
				_, err := utils.GetValidCurrency(mockCurrency)
				Expect(err).To(MatchError("invalid currency"))
			})
		})

		// When("GetValidDueDate recieves a future date ", func() {
		// 	It("returns date", func() {
		//		mockDate := "2032-01-01"
		//		result, err := utils.GetValidDueDate(mockDate)
		// 	})
		//})
		// When("GetValidDueDate recieves a past date ", func() {
		// 	It("throws error", func() {
		//		mockDate := "2000-01-01"
		//		_, err := utils.GetValidDueDate(mockDate)
		// 	})
		// })

	})
})
