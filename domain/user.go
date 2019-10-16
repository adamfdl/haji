package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"golang.org/x/oauth2"
)

type User struct {
	ID             int           `json:"id" db:"id"`
	Email          string        `json:"email" db:"email"`
	FcmID          string        `json:"fcm_id" db:"fcm_id"`
	SeniorityLevel string        `json:"seniority_level" db:"seniority_level"`
	StoryVelocity  storyVelocity `json:"story_velocity" db:"story_velocity"`
	HarvestToken   oauth2.Token  `json:"harvest_token"`
	JiraToken      oauth2.Token  `json:"jira_token"`
}

type storyVelocity map[string]interface{}

func (v storyVelocity) Value() (driver.Value, error) {
	return json.Marshal(v)
}

func (v *storyVelocity) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &v)
}
