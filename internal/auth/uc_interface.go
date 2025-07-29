package auth

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

//go:generate mockery --name UseCase
type UseCase interface {
	Login(ctx context.Context, sc models.Scope, ip LoginInput) (LoginOutput, error)
	RefreshToken(ctx context.Context, sc models.Scope, ip RefreshTokenInput) (RefreshTokenOutput, error)
	Logout(ctx context.Context, sc models.Scope) error
}
