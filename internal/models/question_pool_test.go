package models_test

import (
	. "github.com/jimmykarily/quizmaker/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QuestionPool", func() {
	var poolPath string

	BeforeEach(func() {
		poolPath = "../../tests/assets/question_pool.yaml"
	})

	Describe("NewQuestionPoolFromFile", func() {
		It("returns a new QuestionPool", func() {
			p, err := NewQuestionPoolFromFile(poolPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(p.Questions)).To(Equal(20))
		})
	})
})
