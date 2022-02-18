package url_parsing_strategies_test

import (
	"net/mail"

	"github.com/DusanKasan/parsemail"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"pocok/src/api/process_email/url_parsing_strategies"
)

var _ = Describe("Szamlazz", func() {
	szamlazz := url_parsing_strategies.Szamlazz{}
	var url string
	var testError error

	BeforeEach(func() {
		httpmock.Activate()
	})

	AfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	When("sender is szamlazz with valid html body", func() {
		invoiceSummaryUrl := "https://www.szamlazz.hu/szamla/fiok/73gygabktiedfhxx92aw4mrdjvznjsty2egx?szfejguid=9ge5c6b93jbkwnknt9mt8jns"
		pdfUrl := "https://www.szamlazz.hu/szamla/?action=szamlapdf&szfej_id=226573011&partguid=73gygabktiedfhxx92aw4mrdjvznjsty2egx"

		BeforeEach(func() {
			httpmock.Reset()

			httpmock.RegisterResponder("GET", invoiceSummaryUrl, httpmock.NewStringResponder(200, `
			<div class="flexBox spaceBetween withPadding10">
				<div>
					<div lang="hu" class="tags">Számlaszám</div>
					ASDF-1
				</div>
				<a lang="hu"
					href="/szamla/?action=szamlapdf&amp;szfej_id=226573011&amp;partguid=73gygabktiedfhxx92aw4mrdjvznjsty2egx"
					onclick="event.stopPropagation(); window.open(this.href,'','toolbar=0,scrollbars=1,location=0,statusbar=0,menubar=0,resizable=1,width=900,height=600,left=100,top=100'); return false;"
					class="view-invoice"
				>
					Megnézem
				</a>
			</div>`))

			url, testError = szamlazz.Parse(&parsemail.Email{
				From: []*mail.Address{
					{
						Address: url_parsing_strategies.SzamlazzAddress,
					},
				},
				HTMLBody: `
				<td align="center" valign="top">
					<a href="` + invoiceSummaryUrl + `" style="background-color:#ff6630;color:#ffffff;display:inline-block;font-family:sans-serif;font-size:14px;font-weight:bold;line-height:42px;text-align:center;text-decoration:none;width:300px;text-transform:uppercase" alt="Letöltöm a számlát" target="_blank">Letöltöm a számlát</a>
				</td>`,
			})
		})

		It("returns url", func() {
			Expect(url).To(Equal(pdfUrl))
		})

		It("does not error", func() {
			Expect(testError).To(BeNil())
		})
	})
})
