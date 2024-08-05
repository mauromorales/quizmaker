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

type QuizTemplate struct {
	Questions []Question `yaml:"questions,omitempty"`
}

func NewQuizTemplateFromFile(filePath string) (QuizTemplate, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return QuizTemplate{}, fmt.Errorf("reading file %s: %w", filePath, err)
	}

	return NewQuizTemplate(string(b))
}

func NewQuizTemplate(template string) (QuizTemplate, error) {
	result := QuizTemplate{}

	if err := yaml.Unmarshal([]byte(template), &result); err != nil {
		return result, fmt.Errorf("unmarshaling template: %w", err)
	}

	return result, nil
}
