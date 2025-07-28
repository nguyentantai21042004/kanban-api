package models

import (
	"time"
)

type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	FullName     string     `json:"full_name,omitempty"`
	PasswordHash string     `json:"password_hash,omitempty"`
	AvatarURL    string     `json:"avatar_url,omitempty"`
	IsActive     bool       `json:"is_active,omitempty"`
	CreatedAt    time.Time  `json:"created_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}
