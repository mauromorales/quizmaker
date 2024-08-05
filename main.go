package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jimmykarily/quizmaker/internal/controllers"
	"github.com/jimmykarily/quizmaker/internal/settings"
	settingspkg "github.com/jimmykarily/quizmaker/internal/settings"
)

func main() {
	router := gin.Default()

	var err error
	var settings settingspkg.Settings

	if settings, err = getSettings(); err != nil {
		fmt.Println("Invalid settings: %s", err.Error())
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

func getSettings() (settings.Settings, error) {
	result := settings.Settings{}
	if result.Host = os.Getenv("QUIZMAKER_HOST"); result.Host == "" {
		return result, errors.New("QUIZMAKER_HOST must be set")
	}

	result.InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	result.WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	result.ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return result, nil
}
