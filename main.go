package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jimmykarily/quizmaker/internal/controllers"
	"github.com/jimmykarily/quizmaker/internal/models"
	settingspkg "github.com/jimmykarily/quizmaker/internal/settings"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var questionPoolFlag, databaseStorageDir string

func init() {
	flag.StringVar(&questionPoolFlag, "question-pool", "", "A pool of questions in yaml format")
	flag.StringVar(&databaseStorageDir, "database-storage-dir", "", "The directory where database resides")
	flag.Parse()
}

func main() {
	router := gin.Default()

	var err error
	var settings settingspkg.Settings

	if settings, err = getSettings(); err != nil {
		fmt.Printf("Invalid settings: %s\n", err.Error())
		os.Exit(1)
	}

	if err := autoMigrate(settings.DB); err != nil {
		fmt.Printf("cannot migrate database: %s\n", err.Error())
		os.Exit(1)
	}

	controllers.Settings = settings
	controllers.SetupRoutes(router, controllers.GetRoutes())

	router.Run()
}

func getSettings() (settingspkg.Settings, error) {
	result := settingspkg.Settings{
		InfoLogger:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		WarningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLogger:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	exDir, err := os.Executable() // The directory of the current executable
	if err != nil {
		return result, err
	}

	if databaseStorageDir == "" {
		databaseStorageDir = filepath.Join(filepath.Dir(exDir))
	} else {
		info, err := os.Stat(databaseStorageDir)
		if err != nil {
			if os.IsNotExist(err) {
				return result, fmt.Errorf("database directory does not exist: %s", databaseStorageDir)
			}
			return result, fmt.Errorf("problem with database directory: %w", err)
		}
		if !info.IsDir() {
			return result, fmt.Errorf("database directory '%s' exists but is not a directory: %w", databaseStorageDir, err)
		}
	}

	dbPath := filepath.Join(databaseStorageDir, "database.sql")
	result.DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return result, fmt.Errorf("opening database: %w", err)
	}

	result.QuestionPoolFile = questionPoolFlag
	if result.QuestionPoolFile == "" {
		result.QuestionPoolFile = filepath.Join(filepath.Dir(exDir), "questions.yaml")
	}
	if _, err := os.Stat(result.QuestionPoolFile); err != nil {
		return result, errors.New("no question pool file found (either specified by flag or questions.yaml next to the binary)")
	}

	result.CookieSecret = os.Getenv("QUIZMAKER_COOKIE_SECRET")
	if result.CookieSecret == "" {
		return result, errors.New("QUIZMAKER_COOKIE_SECRET needs to be set to a secret value")
	}

	return result, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Session{})
}
