package models_test

import (
	. "github.com/jimmykarily/quizmaker/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QuizTemplate", func() {
	Describe("NewQuizTemplateFromFile", func() {
		var templatePath string

		BeforeEach(func() {
			templatePath = "../../tests/assets/quiz_template.yaml"
		})

		It("returns a new QuizTemplate", func() {
			t, err := NewQuizTemplateFromFile(templatePath)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(t.Questions)).To(Equal(2))
		})
	})
})
