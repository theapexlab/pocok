package create_email_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "pocok/src/consumers/email_sender/create_email"
)

var _ = Describe("CreateEmail", func() {
	When("it gets the html summary", func() {
		testUrl := "test_api_url"
		testLogoUrl := "https://github.com/theapexlab/pocok/raw/master/assets/pocok-logo.png"
		emailContent, err := GetHtmlSummary(testUrl, testLogoUrl)

		It("returns nil for error", func() {
			Expect(err).To(BeNil())
		})

		It("returns string containing api url", func() {
			Expect(emailContent).To(ContainSubstring(testUrl))
		})
	})
})
