package open_ai_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"

	"pocok/src/mocks/typless/parse_mock_json"
	"pocok/src/services/open_ai"
)

var _ = Describe("OpenAi", func() {
	It("Works", func() {
		extractedData := parse_mock_json.Parse("2021-000022.json")
		text, _ := open_ai.GetPredictedInvoiceData(extractedData.TextBlocks)
		fmt.Println(text)
	})
})
