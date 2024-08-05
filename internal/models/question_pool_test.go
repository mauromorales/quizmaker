package models_test

import (
	. "github.com/jimmykarily/quizmaker/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QuestionPool", func() {
	Describe("NewQuestionPoolFromFile", func() {
		var templatePath string

		BeforeEach(func() {
			templatePath = "../../tests/assets/question_pool.yaml"
		})

		It("returns a new QuestionPool", func() {
			t, err := NewQuestionPoolFromFile(templatePath)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(t.Questions)).To(Equal(2))
		})
	})
})
