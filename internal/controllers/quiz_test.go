package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jimmykarily/quizmaker/internal/controllers"
	"github.com/jimmykarily/quizmaker/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm/clause"
)

var _ = Describe("QuizController test", func() {
	var router *gin.Engine
	var err error
	var w *httptest.ResponseRecorder

	BeforeEach(func() {
		router = gin.Default()
		controllers.SetupRoutes(router, controllers.GetRoutes())

		w = httptest.NewRecorder()

		controllers.Settings.QuestionPoolFile =
			filepath.Join(currentDir, "tests/assets/question_pool.yaml")
	})

	Describe("#New", func() {
		var route controllers.Route

		BeforeEach(func() {
			route, err = controllers.RouteByName("QuizNew")
			Expect(err).ToNot(HaveOccurred())
		})

		It("shows a form for a new quiz", func() {
			req, err := http.NewRequest("GET", route.Path, strings.NewReader("test"))
			Expect(err).ToNot(HaveOccurred())

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK), w.Body.String())
			Expect(w.Body.String()).To(MatchRegexp("form.*action="))
			Expect(w.Body.String()).To(MatchRegexp("Start Quiz"))
		})
	})

	Describe("#Create", func() {
		var email string
		var session models.Session
		var cookie *http.Cookie

		BeforeEach(func() {
			email = "john.doe@example.com"
			w, cookie = performQuizCreateRequest(router, email, nil)
		})

		When("a quiz doesn't exist", func() {
			It("creates a new quiz", func() {
				err := controllers.Settings.DB.Preload(clause.Associations).First(&session).Error
				Expect(err).ToNot(HaveOccurred())

				Expect(w.Body.String()).To(MatchRegexp("Question.*with difficulty"))
				Expect(len(session.Questions)).To(Equal(15))
			})
		})

		When("a quiz already exists", func() {
			BeforeEach(func() {
				err := controllers.Settings.DB.Preload(clause.Associations).First(&session).Error
				Expect(err).ToNot(HaveOccurred())

				Expect(w.Body.String()).To(MatchRegexp("Question.*with difficulty"))
				Expect(len(session.Questions)).To(Equal(15))

				w, cookie = performQuizCreateRequest(router, email, cookie)
			})

			It("doesn't create a new quiz", func() {
				err := controllers.Settings.DB.Preload(clause.Associations).First(&session).Error
				Expect(err).ToNot(HaveOccurred())

				Expect(w.Body.String()).To(MatchRegexp("Question.*with difficulty"))
				Expect(len(session.Questions)).To(Equal(15))
			})
		})
	})
})

func performQuizCreateRequest(router *gin.Engine, email string, cookie *http.Cookie) (*httptest.ResponseRecorder, *http.Cookie) {
	params := map[string]string{
		"email": email,
	}

	route, err := controllers.RouteByName("QuizCreate")
	Expect(err).ToNot(HaveOccurred())

	return performPostWithParams(router, "POST", route.Path, params, cookie)
}
