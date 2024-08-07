package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type QuizOptions struct {
	TotalQuestions     int
	MinDifficulty      int
	MaxDifficulty      int
	QuestionTimeoutSec int
	Questions          QuestionList
}

// Quiz is the collection of questions from a QuestionPool based on QuizOptions.
// It's an intermidiate model used to prepare the Questions that will be stored
// in the database.
type Quiz struct {
	Questions QuestionList `yaml:"questions,omitempty"`
}

func NewQuizWithOpts(opts QuizOptions) (Quiz, error) {
	result := Quiz{}

	result.Questions = opts.Questions.InDifficultyRange(opts.MinDifficulty, opts.MaxDifficulty)

	if opts.TotalQuestions > len(result.Questions) {
		return result, errors.New("not enough questions")
	}

	result.Questions = result.Questions.Limit(opts.TotalQuestions).OrderedByDifficulty()

	return result, nil
}

func (quiz Quiz) PersistForSessionEmail(db *gorm.DB, email string) error {
	s, err := SessionForEmail(db, email)
	if err != nil {
		return fmt.Errorf("looking up session for email %s: %w", email, err)
	}

	return db.Model(&s).Association("Questions").Append(quiz.Questions)
}
