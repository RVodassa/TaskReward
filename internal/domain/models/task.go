package models

import "time"

type Task struct {
	ID          uint       `json:"id"`
	Status      string     `json:"status,omitempty"` // "не завершено", "завершено"
	Description string     `json:"description"`
	Bonus       uint       `json:"bonus"`
	UserID      uint       `json:"user_id,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

func NewTask(description string, bonus uint) *Task {
	now := time.Now().UTC()
	return &Task{
		Description: description,
		Bonus:       bonus,
		CreatedAt:   &now,
	}
}
