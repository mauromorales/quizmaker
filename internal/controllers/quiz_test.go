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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		BeforeEach(func() {
			route, err = controllers.RouteByName("QuizCreate")
			Expect(err).ToNot(HaveOccurred())
		})

		It("creates a new quiz", func() {
			form := url.Values{}
			form.Add("email", "john.doe@example.com")
			encodedForm := form.Encode()

			req, err := http.NewRequest("POST", route.Path, strings.NewReader(encodedForm))
			Expect(err).ToNot(HaveOccurred())

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusFound), w.Body.String())

			redirectURL := w.Result().Header["Location"]
			req, err = http.NewRequest("GET", redirectURL[0], strings.NewReader(encodedForm))
			req.AddCookie(w.Result().Cookies()[0])
			Expect(err).ToNot(HaveOccurred())
			router.ServeHTTP(w, req)

			Expect(w.Body.String()).To(MatchRegexp("Question.*with difficulty"))
		})
	})
})
