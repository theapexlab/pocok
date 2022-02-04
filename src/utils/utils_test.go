package utils_test

import (
	"fmt"
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
		var testError error

		When("it takes correct params", func() {
			BeforeEach(func() {
				mapData = map[string]string{"a": "cica"}
				structData = structType{}
				testError = utils.MapToStruct(mapData, &structData)
			})

			It("does not error", func() {
				Expect(testError).To(BeNil())
			})
			It("the struct has the map data", func() {
				Expect(structData.A).To(Equal("cica"))
			})

		})
	})

	Describe("MapUpdateDataToInvoice", func() {
		When("it receives valid data map as an input", func() {
			It("returns invoice entity with nested services", func() {
				data := map[string]string{
					"invoiceNumber":          "TEST-001",
					"vendorName":             "Vendor Name",
					"vendorEmail":            "test@email.com",
					"accountNumber":          "11736006-20410614",
					"iban":                   "HU19-123456789101112130200005",
					"grossPrice":             "100",
					"vatAmount":              "0",
					"vatRate":                "AAM",
					"currency":               "EUR",
					"service_name_0":         "Service name",
					"service_amount_0":       "db",
					"service_unit_0":         "1",
					"service_unitNetprice_0": "100",
					"service_netPrice_0":     "100",
					"service_grossPrice_0":   "100",
					"service_vatRate_0":      "27%",
					"service_vatAmount_0":    "0",
					"service_currency_0":     "EUR",
					"service_name_1":         "Service name",
					"service_amount_1":       "db",
					"service_unit_1":         "1",
					"service_unitNetprice_1": "100",
					"service_netPrice_1":     "100",
					"service_grossPrice_1":   "100",
					"service_vatRate_1":      "27%",
					"service_vatAmount_1":    "0",
					"service_currency_1":     "EUR",
				}

				invoice, err := utils.MapUpdateDataToInvoice(data)

				Expect(err).To(BeNil())

				for _, service := range invoice.Services {
					Expect(service.Name).To(Equal("Service name"))
					Expect(service.Amount).To(Equal("db"))
					Expect(service.Unit).To(Equal("1"))
					Expect(service.UnitNetPrice).To(Equal("100"))
					Expect(service.NetPrice).To(Equal("100"))
					Expect(service.GrossPrice).To(Equal("100"))
					Expect(service.VatAmount).To(Equal("0"))
					Expect(service.VatRate).To(Equal("27%"))
					Expect(service.Currency).To(Equal("EUR"))
				}

			})
			It("returns invoice entity with empty service array", func() {
				data := map[string]string{
					"invoiceNumber": "TEST-001",
					"vendorName":    "Vendor Name",
					"vendorEmail":   "test@email.com",
					"accountNumber": "11736006-20410614",
					"iban":          "HU19-123456789101112130200005",
					"grossPrice":    "100",
					"vatAmount":     "0",
					"vatRate":       "AAM",
					"currency":      "EUR",
				}

				invoice, err := utils.MapUpdateDataToInvoice(data)

				Expect(err).To(BeNil())
				Expect(len(invoice.Services)).To(Equal(0))

			})
		})
	})

	Describe("Validation test", func() {
		When("GetValidAccountNumber receives a valid 16 or 24 character long account number", func() {
			It("returns Bank account number and recieves no error", func() {
				testBankAccNumbers := []string{
					"12345678-12345678-12345678",
					"12345678 12345678 12345678",
					"12345678 12345678-12345678",
					"123456781234567812345678",
					"12345678-12345678",
					"12345678 12345678",
					"12345678 adfads 12345678",
					"1234567812345678",
				}
				validBankAccNumbers := []string{
					"123456781234567812345678",
					"123456781234567812345678",
					"123456781234567812345678",
					"123456781234567812345678",
					"1234567812345678",
					"1234567812345678",
					"1234567812345678",
					"1234567812345678",
				}
				for i, accNr := range testBankAccNumbers {
					result, err := utils.GetValidAccountNumber(accNr)
					Expect(err).To(BeNil())
					Expect(result).To(Equal(validBankAccNumbers[i]))
				}
			})
		})

		When("GetValidAccountNumber recieves an invalid account number", func() {
			It("returns an error", func() {
				invalidAcountNumbers := []string{
					"12345678-12345678-1234567a",
					"12345678 1  12345678 1234567",
					"12345678-1234567a",
					"12345678-1-12345678",
					"HU35123456781234567",
				}
				for _, accNr := range invalidAcountNumbers {
					_, err := utils.GetValidAccountNumber(accNr)
					Expect(err).To(MatchError("invalid account number"))
				}
			})
		})

		When("GetValidIban recieves valid inputs", func() {
			It("returns an iban code", func() {
				validIbans := []string{
					"HU93116000060000000012345676",
					"HU69119800810030005009212644",
				}
				for i, iban := range validIbans {
					ibanCode, err := utils.GetValidIban(iban)
					Expect(err).To(BeNil())
					Expect(ibanCode).To(Equal(validIbans[i]))
				}
			})
		})
		When("GetValidIban recieves invalid inputs", func() {
			It("returns an error", func() {
				invalid := []string{
					"PL69119800810030005009212644",
					"HU19-123456789101112130200005",
					"12345678-12345678-12345678",
					"12345678-12345678",
				}
				for _, iban := range invalid {
					_, err := utils.GetValidIban(iban)
					Expect(err).ToNot(BeNil())
				}
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

		When("GetValidDueDate recieves a future date ", func() {
			It("returns date", func() {
				mockDate := "2032-01-01"
				result, err := utils.GetValidDueDate(mockDate)
				Expect(result).To(Equal(mockDate))
				Expect(err).To(BeNil())
			})
		})

		When("GetValidPrice recieves a valid price ", func() {
			It("return same price without error", func() {
				testPrices := []string{
					"500",
					"549500.0000",
					"25,000,000",
					"322,50",
					"120.55",
				}
				for i, price := range testPrices {
					result, err := utils.GetValidPrice(price)
					Expect(err).To(BeNil())
					Expect(result).To(Equal(testPrices[i]))

				}
			})
		})

		When("GetValidPrice recieves an invalid price ", func() {
			It("throws error", func() {
				mockPrice := "-500"
				price, err := utils.GetValidPrice(mockPrice)
				fmt.Println("GetValidPrice: " + price)
				Expect(err).To(MatchError("invalid price"))
			})
		})

	})
})
