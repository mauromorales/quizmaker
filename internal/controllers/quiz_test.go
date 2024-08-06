package controllers_test

import (
	"net/http"
	"net/http/httptest"
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

	BeforeEach(func() {
		router = gin.Default()
		controllers.SetupRoutes(router, controllers.GetRoutes())
		route, err = controllers.RouteByName("QuizNew")
		Expect(err).ToNot(HaveOccurred())

		w = httptest.NewRecorder()

		// Change working directory because ginkgo recursively changes this to
		// be the directory of the test. This results in view templates not being
		// found in `Render` function.
		err := os.Chdir(filepath.Join("..", ".."))
		Expect(err).ToNot(HaveOccurred())
	})

	It("returns a new quiz", func() {
		req, err := http.NewRequest("GET", route.Path, strings.NewReader("test"))
		Expect(err).ToNot(HaveOccurred())

		router.ServeHTTP(w, req)

		Expect(w.Code).To(Equal(http.StatusOK), w.Body.String())
		Expect(w.Body.String()).To(MatchRegexp("The quiz questions here"))
	})
})
