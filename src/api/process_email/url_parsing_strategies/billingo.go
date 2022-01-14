package url_parsing_strategies

import (
	"errors"
	"net/http"
	"pocok/src/utils/models"
	"regexp"
	"strings"
)

const BillingoAddress = "noreply@billingo.hu"

type Billingo struct{}

func (b *Billingo) Parse(jsonBody *models.EmailWebhookBody) (string, error) {
	awsTrackUrl, err := parseAwsTrackUrl(jsonBody.Html)
	if err != nil {
		return "", errors.New("can't parse aws track url")
	}

	invoiceSummaryUrl, err := getFinalRedirectUrl(awsTrackUrl)
	if err != nil {
		return "", errors.New("can't get billingo url from aws track url")
	}

	pdfUrl := getPdfUrl(invoiceSummaryUrl)

	return pdfUrl, nil
}

func parseAwsTrackUrl(html string) (string, error) {
	r, err := regexp.Compile(`title="SZÁMLA LETÖLTÉSE" href="(.*)" style`)
	if err != nil {
		return "", err
	}
	matches := r.FindStringSubmatch(html)

	if len(matches) != 2 {
		return "", errors.New("no matches found")
	}

	return matches[1], nil
}

// returns the final url after a serious of redirects
func getFinalRedirectUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	finalURL := resp.Request.URL.String()

	return finalURL, nil
}

/**
the pdf url can be constructed from the invoice summary url
example invoice summary url: https://app.billingo.hu/document-access/default/K90RVdAvQ7gNoq62XvLWJeXq2lDny6aO
example pdf url: https://app.billingo.hu/document-access/K90RVdAvQ7gNoq62XvLWJeXq2lDny6aO/download
*/
func getPdfUrl(invoiceSummaryUrl string) string {
	return strings.Replace(invoiceSummaryUrl, "default/", "", 1) + "/download"
}
