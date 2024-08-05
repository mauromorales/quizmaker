package models_test

import (
	"os"

	. "github.com/jimmykarily/quizmaker/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QuizTemplate", func() {
	Describe("NewQuizTemplate", func() {
		var templatePath string

		BeforeEach(func() {
			b, err := os.ReadFile("../../tests/assets/quiz_template.yaml")
			Expect(err).ToNot(HaveOccurred())
			templatePath = string(b)
		})

		It("returns a new QuizTemplate", func() {
			_, err := NewQuizTemplate(templatePath)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
