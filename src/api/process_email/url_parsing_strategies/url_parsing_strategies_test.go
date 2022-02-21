package url_parsing_strategies_test

import (
	"net/mail"

	"github.com/DusanKasan/parsemail"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"pocok/src/api/process_email/url_parsing_strategies"
)

var _ = Describe("GetPdfUrlFromEmail", func() {
	var url string
	var testError error

	When("there are no strategy for a sender", func() {
		BeforeEach(func() {
			url, testError = url_parsing_strategies.GetPdfUrlFromEmail(&parsemail.Email{
				From: []*mail.Address{
					{
						Address: "test@test.com",
					},
				},
			})
		})

		It("returns nil", func() {
			Expect(url).To(BeEmpty())
		})

		It("errors", func() {
			Expect(testError).To(MatchError(url_parsing_strategies.ErrNoUrlParsingStrategyFound))
		})
	})
})
