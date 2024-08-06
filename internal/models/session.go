package models

import (
	"errors"
	"regexp"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Email string
}

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func NewSessionForEmail(db *gorm.DB, email string) (Session, error) {
	session := Session{Email: email}

	if !ValidEmail(email) {
		return session, errors.New("invalid email")
	}

	result := db.Create(&session)
	if err := result.Error; err != nil {
		return session, err
	}

	return session, nil
}

func ValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
