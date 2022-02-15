package url_parsing_strategies

import (
	"errors"
	HTML "html"
	"io/ioutil"
	"net/http"
	"pocok/src/utils"
	"pocok/src/utils/models"
	"regexp"
	"strings"
)

const SzamlazzAddress = "@szamlazz.hu"

type Szamlazz struct{}

func (sz *Szamlazz) Parse(jsonBody *models.EmailWebhookBody) (string, error) {
	invoiceSummaryUrl, parseErr := sz.parseInvoiceSummaryUrl(jsonBody.Html)
	if parseErr != nil {
		return "", parseErr
	}

	invoiceSummaryHtml, summaryHtmlErr := sz.getInvoiceSummaryHtml(invoiceSummaryUrl)
	if summaryHtmlErr != nil {
		return "", summaryHtmlErr
	}

	pdfUrl, parsePdfUrlError := sz.parsePdfUrl(invoiceSummaryHtml)
	if parsePdfUrlError != nil {
		return "", parsePdfUrlError
	}

	return pdfUrl, nil
}

func (sz *Szamlazz) parseInvoiceSummaryUrl(html string) (string, error) {
	r, regexpError := regexp.Compile(`<a href="(.*)" style=.*>Letöltöm a számlát</a>`)
	if regexpError != nil {
		return "", regexpError
	}
	matches := r.FindStringSubmatch(html)

	if len(matches) != 2 {
		return "", errors.New("no matches found")
	}

	return matches[1], nil
}

func (sz *Szamlazz) getInvoiceSummaryHtml(invoiceSummaryUrl string) (string, error) {
	resp, httpGetError := http.Get(invoiceSummaryUrl)
	if httpGetError != nil {
		utils.LogError("", httpGetError)
		return "", httpGetError
	}

	defer resp.Body.Close()

	html, readAllErr := ioutil.ReadAll(resp.Body)
	if readAllErr != nil {
		utils.LogError("", readAllErr)
		return "", readAllErr
	}

	return string(html), nil
}

func (sz *Szamlazz) parsePdfUrl(html string) (string, error) {
	r, regexErr := regexp.Compile(`(?s)<a.*lang="hu".*href="(.*)".*onclick.*>.*Megnézem.*</a>`)

	if regexErr != nil {
		return "", regexErr
	}
	matches := r.FindStringSubmatch(html)

	if len(matches) != 2 {
		return "", errors.New("no matches found")
	}

	pdfUrl := HTML.UnescapeString(matches[1])

	if !strings.Contains(pdfUrl, "szamlazz.hu") {
		pdfUrl = "https://www.szamlazz.hu" + pdfUrl
	}

	return pdfUrl, nil
}
