package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jimmykarily/quizmaker/internal/controllers"
	"github.com/jimmykarily/quizmaker/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QuestionController test", func() {
	var router *gin.Engine
	var err error
	var w *httptest.ResponseRecorder

	BeforeEach(func() {
		router = gin.Default()
		controllers.SetupRoutes(router, controllers.GetRoutes())

		w = httptest.NewRecorder()
	})

	Describe("#Answer", func() {
		var question models.Question

		When("there is no active session", func() {
			BeforeEach(func() {
				question = models.Question{
					Text: "some question",
				}
				err = controllers.Settings.DB.Save(&question).Error
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns a 401 error", func() {
				params := map[string]string{
					"id":     strconv.Itoa(int(question.ID)),
					"answer": "2",
				}

				path, err := controllers.GetRoutePath("QuestionAnswer",
					map[string]string{"id": strconv.Itoa(int(question.ID))})
				Expect(err).ToNot(HaveOccurred())

				w, _ = performPostWithParams(router, "POST", path, params, nil)
				Expect(w.Code).To(Equal(401))

				err = controllers.Settings.DB.Find(&question).Error
				Expect(err).ToNot(HaveOccurred())

				Expect(question.UserAnswer).To(Equal(0)) // still no answered
			})
		})

		When("there is an active session", func() {
			var cookie *http.Cookie
			var err error
			var session models.Session

			BeforeEach(func() {
				email := "john.doe@example.com"
				cookie, err = controllers.CreateCookie(email, "Firefox")
				Expect(err).ToNot(HaveOccurred())

				session = models.Session{Email: email}
				Expect(controllers.Settings.DB.Create(&session).Error).ToNot(HaveOccurred())
			})

			When("question is not expired and not answered", func() {
				BeforeEach(func() {
					question = models.Question{
						Text:      "some question",
						StartedAt: time.Now().Add(1 * time.Hour),
					}
					err = controllers.Settings.DB.Save(&question).Error
					Expect(err).ToNot(HaveOccurred())
				})

				It("allows answering it", func() {
					params := map[string]string{
						"id":     strconv.Itoa(int(question.ID)),
						"answer": "2",
					}

					path, err := controllers.GetRoutePath("QuestionAnswer",
						map[string]string{"id": strconv.Itoa(int(question.ID))})
					Expect(err).ToNot(HaveOccurred())

					w, _ := performPostWithParams(router, "POST", path, params, cookie)
					Expect(w.Code).To(Equal(http.StatusFound))

					err = controllers.Settings.DB.Find(&question).Error
					Expect(err).ToNot(HaveOccurred())

					Expect(question.UserAnswer).To(Equal(2))
				})
			})
		})
	})
})
