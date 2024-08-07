package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// https://raaaaaaaay86.medium.com/how-to-store-plain-string-slice-by-using-gorm-f855602013e6
type (
	Answers      []string
	QuestionType string
)

type Question struct {
	gorm.Model
	SessionEmail   string
	Session        Session      `gorm:"foreignKey:SessionEmail"`
	Text           string       `yaml:"text,omitempty"`
	Difficulty     int          `yaml:"difficulty,omitempty"`
	Type           QuestionType `yaml:"type,omitempty"`
	RightAnswer    int          `yaml:"rightAnswer,omitempty"`
	UserAnswer     int
	Answers        Answers `yaml:"answers,omitempty" gorm:"type:VARCHAR(255)"`
	AllowedSeconds int     `yaml:"allowedSeconds,omitempty"`
	StartedAt      time.Time
}

func (q Question) Expired() bool {
	isStarted := q.StartedAt.IsZero()
	outOfTime := int(time.Since(q.StartedAt).Seconds()) > q.AllowedSeconds+ALLOWED_SECONDS_SLACK
	notAnswered := (q.UserAnswer == 0)

	return !isStarted && outOfTime && notAnswered
}

// Scan scan value into Jsonb, implements sql.Scanner interface
// https://raaaaaaaay86.medium.com/how-to-store-plain-string-slice-by-using-gorm-f855602013e6
// https://gorm.io/docs/data_types.html
func (a *Answers) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := Answers{}
	err := json.Unmarshal(bytes, &result)
	*a = Answers(result)
	return err
}

// Value return json value, implement driver.Valuer interface
// https://raaaaaaaay86.medium.com/how-to-store-plain-string-slice-by-using-gorm-f855602013e6
// https://gorm.io/docs/data_types.html
func (a Answers) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(Answers(a))
}
