package models_test

import (
	"time"

	. "github.com/jimmykarily/quizmaker/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Session", func() {
	var session Session
	BeforeEach(func() {
		session = Session{}
	})

	Describe("#HasExpiredQuestions", func() {
		When("there are expired questions", func() {
			BeforeEach(func() {
				session.Questions = []Question{
					{
						Text:           "started and expired question",
						StartedAt:      time.Now().Add(-10 * time.Second),
						AllowedSeconds: 5,
					},
					{
						Text:           "not started question",
						AllowedSeconds: 5,
					},
					{
						Text:           "started but not expired question",
						StartedAt:      time.Now(),
						AllowedSeconds: 1000000,
					},
				}
			})

			It("returns true", func() {
				Expect(session.HasExpiredQuestions()).To(BeTrue())
			})
		})

		// An answered question is never considered "expired"
		When("there are no expired questions", func() {
			BeforeEach(func() {
				session.Questions = []Question{
					{
						Text:           "started, anwswered (and 'expired') question",
						StartedAt:      time.Now().Add(-10),
						AllowedSeconds: 5,
						UserAnswer:     2,
					},
					{
						Text:           "not started question",
						AllowedSeconds: 5,
					},
					{
						Text:           "started but not expired question",
						StartedAt:      time.Now(),
						AllowedSeconds: 1000000,
					},
				}
			})

			It("returns false", func() {
				Expect(session.HasExpiredQuestions()).To(BeFalse())
			})
		})
	})

	Describe("#CurrentQuestion", func() {
	})
})
