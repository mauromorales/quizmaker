package models_test

import (
	. "github.com/jimmykarily/quizmaker/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm/clause"
)

var _ = Describe("Quiz", func() {
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
			AvailableQuestions: pool.Questions,
		}
	})

	Describe("#NewQuizWithOpts", func() {
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

	Describe("#PersistForSessionEmail", func() {
		var quiz Quiz
		var email string
		var err error
		var session Session

		BeforeEach(func() {
			quiz, err = NewQuizWithOpts(opts)
			Expect(err).ToNot(HaveOccurred())

			email = "john.doe@example.com"

			session := Session{Email: email}
			Expect(db.Create(&session).Error).ToNot(HaveOccurred())

			Expect(quiz.PersistForSessionEmail(db, email)).ToNot(HaveOccurred())
		})

		It("creates the questions on the database", func() {
			var count int64
			r := db.Model(&Question{}).Count(&count)
			Expect(r.Error).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(4)))
		})

		It("assigns the questions to the specified email/Session", func() {
			Expect(db.Preload(clause.Associations).Find(&session).Error).ToNot(HaveOccurred())
			Expect(len(session.Questions)).To(Equal(4))
		})
	})
})
