package currency_test

import (
	"pocok/src/utils/currency"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetValueFromPrice", func() {
	When("receives input with price amount and currency symbol", func() {
		It("return price as a string value", func() {
			Expect(currency.GetValueFromPrice("7000 Ft")).To(Equal("7000"))
			Expect(currency.GetValueFromPrice("7000 HUF")).To(Equal("7000"))
			Expect(currency.GetValueFromPrice("25,000.00 HUF")).To(Equal("25,000.00"))
			Expect(currency.GetValueFromPrice("25,000,000 HUF")).To(Equal("25,000,000"))
			Expect(currency.GetValueFromPrice("$120")).To(Equal("120"))
			Expect(currency.GetValueFromPrice("120.55 $")).To(Equal("120.55"))
			Expect(currency.GetValueFromPrice("120 USD")).To(Equal("120"))
			Expect(currency.GetValueFromPrice("€300")).To(Equal("300"))
			Expect(currency.GetValueFromPrice("300 €")).To(Equal("300"))
			Expect(currency.GetValueFromPrice("300 EUR")).To(Equal("300"))
			Expect(currency.GetValueFromPrice("€322,50")).To(Equal("322,50"))
			Expect(currency.GetValueFromPrice("442 750 Ft")).To(Equal("442750"))
		})
	})

	When("receives input without valid price amount ", func() {
		It("returns empty string", func() {
			Expect(currency.GetValueFromPrice("EUR")).To(Equal(""))
			Expect(currency.GetValueFromPrice(", USD")).To(Equal(""))
			Expect(currency.GetValueFromPrice(",00 USD")).To(Equal(""))
			Expect(currency.GetValueFromPrice(",0.0 HUF")).To(Equal(""))
		})
	})
})

var _ = Describe("GetCurrencyFromPrice", func() {
	When("recieves string with currency", func() {
		It("returns currency type", func() {
			Expect(currency.GetCurrencyFromString("$150")).To(Equal("USD"))
			Expect(currency.GetCurrencyFromString("400$")).To(Equal("USD"))
			Expect(currency.GetCurrencyFromString("600 $")).To(Equal("USD"))
			Expect(currency.GetCurrencyFromString("630 EUR")).To(Equal("EUR"))
			Expect(currency.GetCurrencyFromString("110€")).To(Equal("EUR"))
			Expect(currency.GetCurrencyFromString("€150")).To(Equal("EUR"))
			Expect(currency.GetCurrencyFromString("28000 HUF")).To(Equal("HUF"))
			Expect(currency.GetCurrencyFromString("15000 FT")).To(Equal("HUF"))
		})

	})

	When("recieves string with no currency symbol", func() {
		It("returns empty string", func() {
			Expect(currency.GetCurrencyFromString("150")).To(Equal(""))
			Expect(currency.GetCurrencyFromString("300 &")).To(Equal(""))
		})
	})
})

var _ = Describe("ConvertPriceToFloat", func() {
	When("recieves price", func() {
		It("returns correct float value", func() {
			Expect(currency.ConvertPriceToFloat("150")).To(Equal(float64(150)))
			Expect(currency.ConvertPriceToFloat("400")).To(Equal(float64(400)))
			Expect(currency.ConvertPriceToFloat("15,000")).To(Equal(float64(15000)))
			Expect(currency.ConvertPriceToFloat("442 750")).To(Equal(float64(442750)))
			Expect(currency.ConvertPriceToFloat("25,000,000")).To(Equal(float64(25000000)))
		})

	})
})
