package models

import (
	"time"

	"github.com/nguyentantai21042004/kanban-api/internal/dbmodels"
)

type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	RoleID       string     `json:"role_id,omitempty"`
	FullName     string     `json:"full_name,omitempty"`
	PasswordHash string     `json:"password_hash,omitempty"`
	AvatarURL    string     `json:"avatar_url,omitempty"`
	IsActive     bool       `json:"is_active,omitempty"`
	CreatedAt    time.Time  `json:"created_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func NewUser(user *dbmodels.User) *User {
	return &User{
		ID:           user.ID,
		Username:     user.Username,
		RoleID:       user.RoleID.String,
		FullName:     user.FullName.String,
		PasswordHash: user.PasswordHash.String,
		AvatarURL:    user.AvatarURL.String,
		IsActive:     user.IsActive.Bool,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		DeletedAt:    &user.DeletedAt.Time,
	}
}

const ADMIN_ROLE = "SUPER_ADMIN"
