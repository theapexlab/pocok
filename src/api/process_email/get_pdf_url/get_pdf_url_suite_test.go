package get_pdf_url_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGetPdfUrl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GetPdfUrl Suite")
}
