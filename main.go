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
	settingspkg "github.com/jimmykarily/quizmaker/internal/settings"
)

var questionPoolFlag string

func init() {
	flag.StringVar(&questionPoolFlag, "question-pool", "", "A pool of questions in yaml format")
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

	controllers.Settings = settings

	setupRoutes(router, controllers.GetRoutes())

	router.Run()
}

func setupRoutes(e *gin.Engine, routes controllers.Routes) {
	e.Static("/assets", "./assets")
	for _, r := range routes {
		e.Handle(r.Method, r.Path, r.Handler)
	}
}

func getSettings() (settingspkg.Settings, error) {
	result := settingspkg.Settings{
		InfoLogger:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		WarningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLogger:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	result.QuestionPoolFile = questionPoolFlag
	if result.QuestionPoolFile == "" {
		ex, err := os.Executable() // The directory of the current executable
		if err != nil {
			return result, err
		}
		result.QuestionPoolFile = filepath.Join(filepath.Dir(ex), "questions.yaml")
	}

	if _, err := os.Stat(result.QuestionPoolFile); err != nil {
		return result, errors.New("no question pool file found (either specified by flag or questions.yaml next to the binary)")
	}

	return result, nil
}
