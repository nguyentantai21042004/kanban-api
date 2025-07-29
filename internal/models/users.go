package models

import (
	"time"
)

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password,omitempty"`
	FullName  string     `json:"full_name,omitempty"`
	AvatarURL string     `json:"avatar_url,omitempty"`
	RoleID    string     `json:"role_id,omitempty"`
	IsActive  bool       `json:"is_active,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
