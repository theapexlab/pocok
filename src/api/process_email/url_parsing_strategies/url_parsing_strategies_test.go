package url_parsing_strategies_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"pocok/src/api/process_email/url_parsing_strategies"
	"pocok/src/utils/models"
)

var _ = Describe("GetPdfUrlFromEmail", func() {
	var url string
	var err error

	When("there are no strategy for a sender", func() {
		BeforeEach(func() {
			url, err = url_parsing_strategies.GetPdfUrlFromEmail(&models.EmailWebhookBody{
				From: []*models.EmailFrom{
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
			Expect(err).To(MatchError(url_parsing_strategies.ErrNoUrlParsingStrategyFound))
		})
	})
})
