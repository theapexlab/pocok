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
})
