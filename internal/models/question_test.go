package models_test

import (
	"time"

	. "github.com/jimmykarily/quizmaker/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Question", func() {
	Describe("#Expired", func() {
		It("returns true when the question is expired", func() {
			question := Question{
				Text:           "started, not anwswered and out of time question",
				StartedAt:      time.Now().Add(-10 * time.Second),
				AllowedSeconds: 5,
			}
			Expect(question.Expired()).To(BeTrue())
		})

		It("returns false when the question is not started", func() {
			question := Question{
				Text:           "not started and not answered question",
				AllowedSeconds: 5,
			}
			Expect(question.Expired()).To(BeFalse())
		})

		It("returns false when the question is still has time", func() {
			question := Question{
				Text:           "started but still going",
				StartedAt:      time.Now(),
				AllowedSeconds: 5000,
			}
			Expect(question.Expired()).To(BeFalse())
		})

		It("returns false when the question is is already answered", func() {
			question := Question{
				Text:           "started, anwswered and out of time question",
				StartedAt:      time.Now().Add(-10 * time.Second),
				AllowedSeconds: 5,
				UserAnswer:     2,
			}
			Expect(question.Expired()).To(BeFalse())
		})
	})
})
