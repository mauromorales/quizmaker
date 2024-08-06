package models_test

import (
	"golang.org/x/exp/maps"

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

	Describe("#InDifficultyRange", func() {
		var pool QuestionPool
		var err error

		BeforeEach(func() {
			pool, err = NewQuestionPoolFromFile(poolPath)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns only questions in the specified difficulty range", func() {
			// sanity check first
			allDifficulties := map[int]int{}
			for _, q := range pool.Questions {
				allDifficulties[q.Difficulty] = 1
			}
			Expect(maps.Keys(allDifficulties)).To(ConsistOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))

			// filtered list now
			questions := pool.Questions.InDifficultyRange(2, 3)
			difficulties := map[int]int{}
			for _, q := range questions {
				difficulties[q.Difficulty] = 1
			}

			Expect(maps.Keys(difficulties)).To(ConsistOf(2, 3))
		})
	})

	Describe("#GenerateQuiz", func() {
		var opts QuizOptions
		var pool QuestionPool
		var err error

		BeforeEach(func() {
			pool, err = NewQuestionPoolFromFile(poolPath)
			Expect(err).ToNot(HaveOccurred())

			opts = QuizOptions{
				TotalQuestions:     4,
				MinDifficulty:      2,
				MaxDifficulty:      4,
				QuestionTimeoutSec: 10,
			}
		})

		It("generates a quiz for the given options", func() {
			q, err := pool.GenerateQuiz(opts)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(q.Questions)).To(Equal(4))
		})

		Describe("validations", func() {
			When("there are not enough questions in the pool", func() {
				BeforeEach(func() {
					opts.TotalQuestions = 100
				})

				It("returns an error", func() {
					_, err := pool.GenerateQuiz(opts)
					Expect(err).To(MatchError("not enough questions"))
				})
			})
		})
	})

	Describe("#Limit", func() {
		var questionPool string

		BeforeEach(func() {
			questionPool = `
questions:
  - text: Q1D1
    difficulty: 1

  - text: Q2D1
    difficulty: 1

  - text: Q1D2
    difficulty: 2

  - text: Q1D3
    difficulty: 3

  - text: Q2D3
    difficulty: 3
`
		})

		It("returns the requested number of questions", func() {
			pool, err := NewQuestionPool(questionPool)
			Expect(err).ToNot(HaveOccurred())
			ql := pool.Questions.Limit(2)
			questions := questionTextFromQuestionList(ql)
			Expect(questions).To(HaveExactElements("Q1D1", "Q1D2"))

			pool, err = NewQuestionPool(questionPool)
			Expect(err).ToNot(HaveOccurred())
			ql = pool.Questions.Limit(3)
			questions = questionTextFromQuestionList(ql)
			Expect(questions).To(HaveExactElements("Q1D1", "Q1D2", "Q1D3"))

			pool, err = NewQuestionPool(questionPool)
			Expect(err).ToNot(HaveOccurred())
			ql = pool.Questions.Limit(4)
			questions = questionTextFromQuestionList(ql)
			Expect(questions).To(HaveExactElements("Q1D1", "Q1D2", "Q1D3", "Q2D1"))

			pool, err = NewQuestionPool(questionPool)
			Expect(err).ToNot(HaveOccurred())
			ql = pool.Questions.Limit(5)
			questions = questionTextFromQuestionList(ql)
			Expect(questions).To(HaveExactElements("Q1D1", "Q1D2", "Q1D3", "Q2D1", "Q2D3"))
		})
	})
})

func questionTextFromQuestionList(ql QuestionList) []string {
	result := []string{}
	for _, q := range ql {
		result = append(result, q.Text)
	}

	return result
}
