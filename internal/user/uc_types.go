package user

import (
	"time"

	"github.com/nguyentantai21042004/kanban-api/internal/models"
)

type CreateInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	RoleID   string `json:"role_id"`
}

type UpdateProfileInput struct {
	FullName  string `json:"full_name" binding:"required"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type UserOutput struct {
	User models.User `json:"user"`
}

type GetUserOutput struct {
	Users []models.User `json:"users"`
}

type GetOneInput struct {
	Username string
}

type ListInput struct {
	Filter Filter
}

type Filter struct {
	IDs []string
}

// Dashboard aggregation for users
type DashboardInput struct {
	From time.Time
	To   time.Time
}

type UsersDashboardOutput struct {
	Total  int64
	Active int64
	Growth float64
}
