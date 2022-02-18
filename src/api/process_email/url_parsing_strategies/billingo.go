package url_parsing_strategies

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/DusanKasan/parsemail"
)

const BillingoAddress = "noreply@billingo.hu"

type Billingo struct{}

func (b *Billingo) Parse(email *parsemail.Email) (string, error) {
	awsTrackUrl, parseError := b.parseAwsTrackUrl(email.HTMLBody)
	if parseError != nil {
		return "", errors.New("can't parse aws track url")
	}

	invoiceSummaryUrl, getRedirectUrlError := b.getFinalRedirectUrl(awsTrackUrl)
	if getRedirectUrlError != nil {
		return "", errors.New("can't get billingo url from aws track url")
	}

	pdfUrl := b.getPdfUrl(invoiceSummaryUrl)

	return pdfUrl, nil
}

func (b *Billingo) parseAwsTrackUrl(html string) (string, error) {
	r, regexpError := regexp.Compile(`title="SZÁMLA LETÖLTÉSE" href="(.*)" style`)
	if regexpError != nil {
		return "", regexpError
	}
	matches := r.FindStringSubmatch(html)

	if len(matches) != 2 {
		return "", errors.New("no matches found")
	}

	return matches[1], nil
}

// returns the final url after a serious of redirects
func (b *Billingo) getFinalRedirectUrl(url string) (string, error) {
	resp, httpGetError := http.Get(url)
	if httpGetError != nil {
		return "", httpGetError
	}

	finalURL := resp.Request.URL.String()

	return finalURL, nil
}

/**
the pdf url can be constructed from the invoice summary url
example invoice summary url: https://app.billingo.hu/document-access/default/K90RVdAvQ7gNoq62XvLWJeXq2lDny6aO
example pdf url: https://app.billingo.hu/document-access/K90RVdAvQ7gNoq62XvLWJeXq2lDny6aO/download
*/
func (b *Billingo) getPdfUrl(invoiceSummaryUrl string) string {
	return strings.Replace(invoiceSummaryUrl, "default/", "", 1) + "/download"
}
