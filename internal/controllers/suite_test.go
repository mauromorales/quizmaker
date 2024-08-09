package controllers_test

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jimmykarily/quizmaker/internal/controllers"
	"github.com/jimmykarily/quizmaker/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controllers Suite")
}

var originalWorkingDir string
var currentDir string

var _ = BeforeEach(func() {
	// reset the db before each test
	testDbPath, err := filepath.Abs(filepath.Join("..", "..", "tests", "database.sql"))
	Expect(err).ToNot(HaveOccurred())
	err = os.RemoveAll(testDbPath)
	Expect(err).ToNot(HaveOccurred())

	controllers.Settings.DB, err = gorm.Open(sqlite.Open(testDbPath), &gorm.Config{})
	Expect(err).ToNot(HaveOccurred())

	err = models.AutoMigrate(controllers.Settings.DB)
	Expect(err).ToNot(HaveOccurred())

	// Change working directory because ginkgo recursively changes this to
	// be the directory of the test. This results in view templates not being
	// found in `Render` function.
	originalWorkingDir, err = os.Getwd()
	Expect(err).ToNot(HaveOccurred())
	err = os.Chdir(filepath.Join("..", ".."))
	Expect(err).ToNot(HaveOccurred())
	currentDir, err = os.Getwd()
	Expect(err).ToNot(HaveOccurred())

	cookieSecret, err := generateSecret()
	Expect(err).ToNot(HaveOccurred())

	controllers.Settings.CookieSecret = cookieSecret
})

var _ = AfterEach(func() {
	Expect(os.Chdir(originalWorkingDir)).ToNot(HaveOccurred())
})

func generateSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func performPostWithParams(router *gin.Engine, verb, path string, params map[string]string, cookie *http.Cookie) (*httptest.ResponseRecorder, *http.Cookie) {
	fmt.Printf("cookie = %+v\n", cookie)
	w := httptest.NewRecorder()
	form := url.Values{}
	for k, v := range params {
		form.Add(k, v)
	}
	encodedForm := form.Encode()

	req, err := http.NewRequest(verb, path, strings.NewReader(encodedForm))
	Expect(err).ToNot(HaveOccurred())

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if cookie != nil {
		req.AddCookie(cookie)
	}

	router.ServeHTTP(w, req)
	var newCookie *http.Cookie
	if cookies := w.Result().Cookies(); len(cookies) > 0 {
		newCookie = w.Result().Cookies()[0]
	}
	// no new cookie has been sent, keep the old one
	if newCookie == nil {
		newCookie = cookie
	}

	// Stop here if there was no redirection
	if w.Code != http.StatusFound {
		return w, newCookie
	}

	redirectURL := w.Result().Header["Location"]
	req, err = http.NewRequest("GET", redirectURL[0], strings.NewReader(encodedForm))
	Expect(err).ToNot(HaveOccurred())
	if newCookie != nil {
		req.AddCookie(newCookie)
	}
	router.ServeHTTP(w, req)

	return w, newCookie
}
