package controllers_test

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"testing"

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

// reset the db before each test
var _ = BeforeEach(func() {
	testDbPath, err := filepath.Abs(filepath.Join("..", "..", "tests", "database.sql"))
	Expect(err).ToNot(HaveOccurred())
	err = os.RemoveAll(testDbPath)
	Expect(err).ToNot(HaveOccurred())

	controllers.Settings.DB, err = gorm.Open(sqlite.Open(testDbPath), &gorm.Config{})
	Expect(err).ToNot(HaveOccurred())

	err = models.AutoMigrate(controllers.Settings.DB)
	Expect(err).ToNot(HaveOccurred())

	cookieSecret, err := generateSecret()
	Expect(err).ToNot(HaveOccurred())

	controllers.Settings.CookieSecret = cookieSecret
})

func generateSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
