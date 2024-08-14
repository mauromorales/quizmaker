package models

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Email     string
	Nickname  string
	Score     int
	Complete  bool
	Questions []Question `gorm:"foreignKey:SessionEmail;references:Email"`
}

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func NewSession(db *gorm.DB, email, nickname string) (Session, error) {
	session := Session{Email: email, Nickname: nickname}

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
	// sort by index
	sort.Slice(s.Questions, func(i, j int) bool {
		return s.Questions[i].Index < s.Questions[j].Index
	})

	// if there is an already started question, return that, no matter the Index
	// otherwise it will expire while answering another one. This should not happen
	// if the questions are presented in order of Index but this code handles this
	// anyway.
	for _, q := range s.Questions {
		if !q.StartedAt.IsZero() && q.UserAnswer == 0 && !q.Expired() {
			return q, nil
		}
	}

	// return the first unanswered question by Index
	for _, q := range s.Questions {
		if q.UserAnswer == 0 && !q.Expired() {
			return q, nil
		}
	}

	// no unanswered question found
	return Question{}, nil
}

// UpdateCacheColumns calculates the current "Score" value based only on
// answered and expired questions. Expired questions with no answer
// are considered "wrong".
// It also calculated the value of the "Completed" column. A session is complete
// when all questions are answered or expired.
func (s *Session) UpdateCacheColumns() {
	correctAnswers := 0
	completeQuestions := 0
	totalQuestions := len(s.Questions)

	for _, q := range s.Questions {
		if q.Expired() {
			completeQuestions++ // consider it a wrong answer
			continue
		}

		// ignore not started or in-progress questions (we already handled expired above)
		if q.StartedAt.IsZero() || q.UserAnswer == 0 {
			continue
		}

		if q.UserAnswer == q.RightAnswer {
			correctAnswers++
		}
		completeQuestions++ // right or wrong, count it in
	}

	if completeQuestions == totalQuestions {
		s.Complete = true
	}
	s.Score = int(math.Round(float64(correctAnswers) / float64(completeQuestions) * 100))
}

// EmailObfuscated obfuscates an email address by replacing characters with dots,
// except for the first and last characters of the username and domain parts.
func (s Session) EmailObfuscated() string {
	parts := strings.Split(s.Email, "@")
	if len(parts) != 2 {
		return s.Email // Return the original email if it's invalid
	}

	obfuscatedUsername := obfuscateString(parts[0])
	obfuscatedDomain := obfuscateDomain(parts[1])
	fmt.Printf("obfuscatedDomain = %+v\n", obfuscatedDomain)
	fmt.Printf("obfuscatedUsername = %+v\n", obfuscatedUsername)

	return obfuscatedUsername + "@" + obfuscatedDomain
}

func obfuscateString(part string) string {
	if len(part) <= 2 {
		return part // Do not obfuscate if the part is too short
	}

	// Keep the first and last characters, replace the rest with dots
	return string(part[0]) + strings.Repeat(".", len(part)-2) + string(part[len(part)-1])
}

// obfuscateDomain obfuscates the domain part of the email while preserving the TLD.
func obfuscateDomain(domain string) string {
	domainParts := strings.Split(domain, ".")
	if len(domainParts) < 2 {
		return obfuscateString(domain) // If the domain is not well-formed, treat it as a simple part
	}

	// Obfuscate everything except the last part (the TLD)
	for i := 0; i < len(domainParts)-1; i++ {
		domainParts[i] = obfuscateString(domainParts[i])
	}

	return strings.Join(domainParts, ".")
}
