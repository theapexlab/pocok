package create_invoice_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCreateInvoice(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateInvoice Suite")
}
