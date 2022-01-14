package parse_expense_analysis_test

import (
	"encoding/json"
	"io/ioutil"
	"pocok/src/consumers/textractor/parse_expense_analysis"

	"github.com/aws/aws-sdk-go-v2/service/textract"
	. "github.com/onsi/ginkgo/v2"
)

func parseMockJson(filename string) *textract.GetExpenseAnalysisOutput {
	var expenseAnalysisOutput *textract.GetExpenseAnalysisOutput

	mock, readFileErr := ioutil.ReadFile("../../../mocks/textract/" + filename)
	if readFileErr != nil {
		panic("Failed to read mock file")
	}

	if err := json.Unmarshal(mock, &expenseAnalysisOutput); err != nil {
		panic("Failed to unmarshal mock file")
	}

	return expenseAnalysisOutput
}

var _ = Describe("ParseExpenseAnalysis", func() {
	var expenseAnalysisOutput *textract.GetExpenseAnalysisOutput

	When("", func() {
		BeforeEach(func() {
			expenseAnalysisOutput = parseMockJson("angelina.json")
		})

		It("works", func() {
			parse_expense_analysis.ParseExpenseAnalysis(expenseAnalysisOutput)
		})
	})
})
