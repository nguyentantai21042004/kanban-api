package http

import (
	"github.com/nguyentantai21042004/kanban-api/internal/user"
)

type createReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

func (req createReq) toInput() user.CreateInput {
	return user.CreateInput{
		Username: req.Username,
		Password: req.Password,
		FullName: req.FullName,
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
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url,omitempty"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (h handler) newItem(o user.UserOutput) userItem {
	return userItem{
		ID:        o.User.ID,
		Username:  o.User.Username,
		FullName:  o.User.FullName,
		AvatarURL: o.User.AvatarURL,
		IsActive:  o.User.IsActive,
		CreatedAt: o.User.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: o.User.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
