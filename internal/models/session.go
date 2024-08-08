package models

import (
	"errors"
	"regexp"
	"sort"

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

// CurrentQuestion returns the first non-expired question
// If the session has an expired question, this method returns an error
// (there is no "CurrentQuestion" since the whole Session is expired)
func (s Session) CurrentQuestion() (Question, error) {
	if s.HasExpiredQuestions() {
		return Question{}, errors.New("time is out, quiz is expired")
	}

	// sort by index
	sort.Slice(s.Questions, func(i, j int) bool {
		return s.Questions[i].Index < s.Questions[j].Index
	})

	// if there is an already started question, return that, no matter the Index
	// otherwise it will expire will answering another one. This should not happen
	// if the questions are presented in order of Index but this code handles this
	// anyway.
	for _, q := range s.Questions {
		if !q.StartedAt.IsZero() && q.UserAnswer == 0 {
			return q, nil
		}
	}

	// return the first unanswered question by Index
	for _, q := range s.Questions {
		if q.UserAnswer == 0 {
			return q, nil
		}
	}

	// no unanswered question found
	return Question{}, nil
}
