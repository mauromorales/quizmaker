package models_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jimmykarily/quizmaker/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

// reset the db before each test
var _ = BeforeEach(func() {
	testDbPath, err := filepath.Abs(filepath.Join("..", "..", "tests", "database.sql"))
	Expect(err).ToNot(HaveOccurred())
	err = os.RemoveAll(testDbPath)
	Expect(err).ToNot(HaveOccurred())

	db, err = gorm.Open(sqlite.Open(testDbPath), &gorm.Config{})
	Expect(err).ToNot(HaveOccurred())

	err = models.AutoMigrate(db)
	Expect(err).ToNot(HaveOccurred())
})
