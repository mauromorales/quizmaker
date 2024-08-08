package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
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
	var route controllers.Route
	var err error
	var w *httptest.ResponseRecorder
	var originalWorkingDir string

	BeforeEach(func() {
		router = gin.Default()
		controllers.SetupRoutes(router, controllers.GetRoutes())

		w = httptest.NewRecorder()

		// Change working directory because ginkgo recursively changes this to
		// be the directory of the test. This results in view templates not being
		// found in `Render` function.
		originalWorkingDir, err = os.Getwd()
		Expect(err).ToNot(HaveOccurred())
		err = os.Chdir(filepath.Join("..", ".."))
		Expect(err).ToNot(HaveOccurred())
		currentDir, err := os.Getwd()
		Expect(err).ToNot(HaveOccurred())

		controllers.Settings.QuestionPoolFile =
			filepath.Join(currentDir, "tests/assets/question_pool.yaml")
	})

	AfterEach(func() {
		err = os.Chdir(originalWorkingDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("#New", func() {
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

			route, err = controllers.RouteByName("QuizCreate")
			Expect(err).ToNot(HaveOccurred())

			w, cookie = performRequest(router, route, email, nil)
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
			})

			It("doesn't create a new quiz", func() {
				w, cookie = performRequest(router, route, email, cookie)
				err := controllers.Settings.DB.Preload(clause.Associations).First(&session).Error
				Expect(err).ToNot(HaveOccurred())

				Expect(w.Body.String()).To(MatchRegexp("Question.*with difficulty"))
				Expect(len(session.Questions)).To(Equal(15))
			})
		})
	})
})

func performRequest(router *gin.Engine, route controllers.Route, email string, cookie *http.Cookie) (*httptest.ResponseRecorder, *http.Cookie) {
	w := httptest.NewRecorder()
	form := url.Values{}
	form.Add("email", email)
	encodedForm := form.Encode()

	req, err := http.NewRequest("POST", route.Path, strings.NewReader(encodedForm))
	Expect(err).ToNot(HaveOccurred())

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if cookie != nil {
		req.AddCookie(cookie)
	}

	router.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusFound), w.Body.String())

	var newCookie *http.Cookie
	if cookies := w.Result().Cookies(); len(cookies) > 0 {
		newCookie = w.Result().Cookies()[0]
	}

	redirectURL := w.Result().Header["Location"]
	req, err = http.NewRequest("GET", redirectURL[0], strings.NewReader(encodedForm))
	Expect(err).ToNot(HaveOccurred())
	if newCookie == nil {
		newCookie = cookie
	}
	req.AddCookie(newCookie)
	router.ServeHTTP(w, req)

	return w, newCookie
}
