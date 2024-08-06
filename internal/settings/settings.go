package settings

import (
	"log"

	"gorm.io/gorm"
)

type Settings struct {
	InfoLogger       *log.Logger
	WarningLogger    *log.Logger
	ErrorLogger      *log.Logger
	QuestionPoolFile string
	DB               *gorm.DB
	CookieSecret     string
}
