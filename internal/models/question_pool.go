package models

import (
	"fmt"
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

type QuestionPool struct {
	Questions []Question `yaml:"questions,omitempty"`
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
