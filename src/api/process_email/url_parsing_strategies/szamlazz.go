package url_parsing_strategies

import (
	"errors"
	HTML "html"
	"pocok/src/utils"
	"regexp"
	"strings"

	"github.com/DusanKasan/parsemail"
)

const SzamlazzAddress = "@szamlazz.hu"

type Szamlazz struct{}

func (sz *Szamlazz) Parse(email *parsemail.Email) (string, error) {
	invoiceSummaryUrl, parseErr := sz.parseInvoiceSummaryUrl(email.HTMLBody)
	if parseErr != nil {
		return "", parseErr
	}

	invoiceSummaryHtmlBytes, summaryHtmlErr := utils.DownloadFile(invoiceSummaryUrl)
	if summaryHtmlErr != nil {
		return "", summaryHtmlErr
	}
	invoiceSummaryHtml := string(invoiceSummaryHtmlBytes)

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
