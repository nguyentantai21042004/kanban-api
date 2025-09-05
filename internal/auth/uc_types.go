package auth

import "github.com/nguyentantai21042004/kanban-api/internal/models"

type LoginInput struct {
	Username string
	Password string
}

type LoginOutput struct {
	AssToken string
	User     models.User
	Role     models.Role
}

type RefreshTokenInput struct {
	RfrToken string
}

type RefreshTokenOutput struct {
	AssToken string
	RfrToken string
}
