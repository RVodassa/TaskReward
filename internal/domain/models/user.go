package models

import (
	"time"
)

type User struct {
	ID           uint       `json:"id"`
	Login        string     `json:"login,omitempty"`
	PasswordHash string     `json:"-"`
	ReferID      uint       `json:"refer_id,omitempty"`
	Balance      uint       `json:"balance"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
}

// NewUser создает новый инстанс пользователя
func NewUser(login string, password string, referID uint) *User {
	now := time.Now().UTC()
	return &User{
		Login:        login,
		PasswordHash: password,
		ReferID:      referID,
		CreatedAt:    &now,
	}
}
