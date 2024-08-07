package models

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type QuestionPool struct {
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
