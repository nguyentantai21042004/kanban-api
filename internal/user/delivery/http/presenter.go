package http

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
)

type createReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	RoleID   string `json:"role_id" binding:"required"`
}

func (req createReq) toInput() user.CreateInput {
	return user.CreateInput{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		RoleID:   req.RoleID,
	}
}

type updateProfileReq struct {
	FullName  string `json:"full_name" binding:"required"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

func (req updateProfileReq) toInput() user.UpdateProfileInput {
	return user.UpdateProfileInput{
		FullName:  req.FullName,
		AvatarURL: req.AvatarURL,
	}
}

type userItem struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url,omitempty"`
	RoleID    string `json:"role_id"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (h handler) newItem(o user.UserOutput) userItem {
	return userItem{
		ID:        o.User.ID,
		Email:     o.User.Email,
		FullName:  o.User.FullName,
		AvatarURL: o.User.AvatarURL,
		RoleID:    o.User.RoleID,
		IsActive:  o.User.IsActive,
		CreatedAt: o.User.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: o.User.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
