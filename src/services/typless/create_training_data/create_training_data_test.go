package create_training_data_test

import (
	"pocok/src/mocks"
	"pocok/src/services/typless"
	"pocok/src/services/typless/create_training_data"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateTrainingData", func() {
	var trainingData *typless.TrainingData

	When("gets valid invoice", func() {
		BeforeEach(func() {
			trainingData = create_training_data.CreateTrainingData(&mocks.MockInvoice)
		})

		It("", func() {
			// fmt.Println(trainingData)
			Expect(len(trainingData.LearningFields)).To(Equal(len(typless.ExtractDataToInvoiceMap)))
		})
	})
})
