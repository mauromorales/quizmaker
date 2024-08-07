package models

import "errors"

type QuizOptions struct {
	TotalQuestions     int
	MinDifficulty      int
	MaxDifficulty      int
	QuestionTimeoutSec int
	Questions          QuestionList
}

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
