package url_parsing_strategies_test

import (
	"net/http"
	"net/mail"

	"github.com/DusanKasan/parsemail"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"pocok/src/api/process_email/url_parsing_strategies"
)

var _ = Describe("Billingo", func() {
	billingo := url_parsing_strategies.Billingo{}
	var url string
	var testError error

	BeforeEach(func() {
		httpmock.Activate()
	})

	AfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	When("sender is billingo with valid html body", func() {
		awsTrackUrl := "https://2g2zw50k.r.eu-central-1.awstrack.me/L0/https:%2F%2Fapp.billingo.hu%2Fdocument-access%2FK90RVdAvQ7gNoq62XvLWJeXq2lDny6aO/1/0107017dfb72c633-ff7877cf-fac2-415e-bb98-6f1709ae7470-000000/OHClr1vWQbzCIWXrIxHWiBDQDHg=30"
		invoiceSummaryUrl := "https://app.billingo.hu/document-access/default/K90RVdAvQ7gNoq62XvLWJeXq2lDny6aO"
		pdfUrl := "https://app.billingo.hu/document-access/K90RVdAvQ7gNoq62XvLWJeXq2lDny6aO/download"

		BeforeEach(func() {
			httpmock.Reset()

			httpmock.RegisterResponder("GET", awsTrackUrl,
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(301, "")
					resp.Header.Add("Location", invoiceSummaryUrl)
					return resp, nil
				})
			httpmock.RegisterResponder("GET", invoiceSummaryUrl, httpmock.NewStringResponder(200, ""))

			url, testError = billingo.Parse(&parsemail.Email{
				From: []*mail.Address{
					{
						Address: url_parsing_strategies.BillingoAddress,
					},
				},
				HTMLBody: `
				<td align="right" valign="middle" style="font-family:Arial,&#39;Helvetica Neue&#39;,Helvetica,sans-serif;font-size:14px;padding:17px;border-collapse:separate!important;border:2px none #707070;border-radius:5px;background-color:#78d230">
					<a title="SZÁMLA LETÖLTÉSE" href="` + awsTrackUrl + `" style="font-weight:bold;letter-spacing:normal;line-height:100%;text-align:center;text-decoration:none;color:#ffffff" target="_blank">SZÁMLA LETÖLTÉSE</a>
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
