package settings

import "log"

type Settings struct {
	Host          string
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
}
