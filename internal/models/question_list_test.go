package models_test

import (
	"golang.org/x/exp/maps"

	. "github.com/jimmykarily/quizmaker/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QuestionList", func() {
	var poolPath string

	BeforeEach(func() {
		poolPath = "../../tests/assets/question_pool.yaml"
	})

	Describe("#Valid", func() {
		var list QuestionList

		BeforeEach(func() {
			pool, err := NewQuestionPool(`
questions:
  - text: Invalid Question (no "rightAnswer" set)
    ID: 1
    difficulty: 3
    type: multiple-choice
    answers:
    - answer 1
    - answer 2

  - text: Invalid Question ("rightAnswer" set to non-existent index)
    ID: 2
    difficulty: 3
    type: multiple-choice
    rightAnswer: 10
    answers:
    - answer 1
    - answer 2

  - text: Valid question
    ID: 3
    difficulty: 1
    type: boolean
    rightAnswer: 2
    answers:
    - True
    - False
`)
			Expect(err).ToNot(HaveOccurred())
			Expect(db.Save(&pool.Questions).Error).ToNot(HaveOccurred())
			list = pool.Questions
		})

		It("returns only questions that are valid", func() {
			// sanity check first
			allIDs := map[int]int{}
			for _, q := range list {
				allIDs[int(q.ID)] = 1
			}
			Expect(maps.Keys(allIDs)).To(ConsistOf(1, 2, 3))

			// filtered list now
			questions := list.Valid()
			ids := map[int]int{}
			for _, q := range questions {
				ids[int(q.ID)] = 1
			}

			Expect(maps.Keys(ids)).To(ConsistOf(3))
		})
	})

	Describe("#InDifficultyRange", func() {
		var list QuestionList

		BeforeEach(func() {
			pool, err := NewQuestionPoolFromFile(poolPath)
			Expect(err).ToNot(HaveOccurred())
			list = pool.Questions
		})

		It("returns only questions in the specified difficulty range", func() {
			// sanity check first
			allDifficulties := map[int]int{}
			for _, q := range list {
				allDifficulties[q.Difficulty] = 1
			}
			Expect(maps.Keys(allDifficulties)).To(ConsistOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))

			// filtered list now
			questions := list.InDifficultyRange(2, 3)
			difficulties := map[int]int{}
			for _, q := range questions {
				difficulties[q.Difficulty] = 1
			}

			Expect(maps.Keys(difficulties)).To(ConsistOf(2, 3))
		})
	})

	Describe("#Limit", func() {
		var list QuestionList

		BeforeEach(func() {
			questionPool := `
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
			pool, err := NewQuestionPool(questionPool)
			Expect(err).ToNot(HaveOccurred())

			list = pool.Questions
		})

		It("returns the requested number of questions", func() {
			ql := list.Limit(2)
			questions := questionTextFromQuestionList(ql)
			Expect(questions).To(HaveExactElements("Q1D1", "Q1D2"))

			ql = list.Limit(3)
			questions = questionTextFromQuestionList(ql)
			Expect(questions).To(HaveExactElements("Q1D1", "Q1D2", "Q1D3"))

			ql = list.Limit(4)
			questions = questionTextFromQuestionList(ql)
			Expect(questions).To(HaveExactElements("Q1D1", "Q1D2", "Q1D3", "Q2D1"))

			ql = list.Limit(5)
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
