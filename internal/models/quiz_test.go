package models_test

import (
	. "github.com/jimmykarily/quizmaker/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Quiz", func() {
	Describe("#NewQuizWithOpts", func() {
		var opts QuizOptions

		BeforeEach(func() {
			poolPath := "../../tests/assets/question_pool.yaml"
			pool, err := NewQuestionPoolFromFile(poolPath)
			Expect(err).ToNot(HaveOccurred())

			opts = QuizOptions{
				TotalQuestions:     4,
				MinDifficulty:      2,
				MaxDifficulty:      4,
				QuestionTimeoutSec: 10,
				Questions:          pool.Questions,
			}
		})

		It("generates a quiz for the given options", func() {
			q, err := NewQuizWithOpts(opts)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(q.Questions)).To(Equal(4))
		})

		Describe("validations", func() {
			When("there are not enough questions in the pool", func() {
				BeforeEach(func() {
					opts.TotalQuestions = 100
				})

				It("returns an error", func() {
					_, err := NewQuizWithOpts(opts)
					Expect(err).To(MatchError("not enough questions"))
				})
			})
		})
	})
})
