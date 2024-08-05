package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jimmykarily/quizmaker/internal/controllers"
	settingspkg "github.com/jimmykarily/quizmaker/internal/settings"
)

func main() {
	router := gin.Default()

	var err error
	var settings settingspkg.Settings

	if settings, err = getSettings(); err != nil {
		fmt.Printf("Invalid settings: %s", err.Error())
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

	return result, nil
}
