package models

import (
	"errors"
	"regexp"

	"gorm.io/gorm"
)

// Rendering web pages might take some time.
// Add some "slack" to the allowed seconds to cover for that.
const ALLOWED_SECONDS_SLACK = 2

type Session struct {
	gorm.Model
	Email     string
	Questions []Question `gorm:"foreignKey:SessionEmail;references:Email"`
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

func SessionForEmail(db *gorm.DB, email string) (Session, error) {
	var session Session
	result := db.First(&session, "email = ?", email)
	if err := result.Error; err != nil {
		return session, err
	}

	return session, nil
}

func ValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// HasExpiredQuestions returns `true` if there is at least one expired Question
func (s Session) HasExpiredQuestions() bool {
	for _, q := range s.Questions {
		if q.Expired() {
			return true
		}
	}
	return false
}

// CurrentQuestion returns the first non-completed question
func (s Session) CurrentQuestion() Question {
	return Question{} // TODO
}
