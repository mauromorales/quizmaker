package models

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"gopkg.in/yaml.v3"
)

type QuestionType string

type Question struct {
	Text        string       `yaml:"text,omitempty"`
	Difficulty  int          `yaml:"difficulty,omitempty"`
	Type        QuestionType `yaml:"type,omitempty"`
	RightAnswer int          `yaml:"rightAnswer,omitempty"`
	Answers     []string     `yaml:"answers,omitempty"`
}

type QuestionList []Question

type QuestionPool struct {
	Questions QuestionList `yaml:"questions,omitempty"`
}

type QuizOptions struct {
	TotalQuestions     int
	MinDifficulty      int
	MaxDifficulty      int
	QuestionTimeoutSec int
}

type Quiz struct {
	Questions QuestionList `yaml:"questions,omitempty"`
}

func NewQuestionPoolFromFile(filePath string) (QuestionPool, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return QuestionPool{}, fmt.Errorf("reading file %s: %w", filePath, err)
	}

	return NewQuestionPool(string(b))
}

func NewQuestionPool(template string) (QuestionPool, error) {
	result := QuestionPool{}

	if err := yaml.Unmarshal([]byte(template), &result); err != nil {
		return result, fmt.Errorf("unmarshaling template: %w", err)
	}

	return result, nil
}

func (ql QuestionList) InDifficultyRange(min, max int) QuestionList {
	result := QuestionList{}
	for _, q := range ql {
		if q.Difficulty <= max && q.Difficulty >= min {
			result = append(result, q)
		}
	}

	return result
}

func (ql QuestionList) Suffled() QuestionList {
	dest := make(QuestionList, len(ql))
	perm := rand.Perm(len(ql))
	for i, v := range perm {
		dest[v] = ql[i]
	}

	return dest
}

func (ql QuestionList) Limit(limit int) QuestionList {
	if len(ql) > limit {
		return ql[0:limit]
	}

	return ql
}

func (qp QuestionPool) GenerateQuiz(opts QuizOptions) (Quiz, error) {
	result := Quiz{}

	result.Questions = qp.Questions.
		InDifficultyRange(opts.MinDifficulty, opts.MaxDifficulty).
		Suffled()

	if opts.TotalQuestions > len(result.Questions) {
		return result, errors.New("not enough questions")
	}

	result.Questions = result.Questions.Limit(opts.TotalQuestions)

	return result, nil
}
