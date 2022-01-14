package parse_expense_analysis_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestParseExpenseAnalysis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ParseExpenseAnalysis Suite")
}
