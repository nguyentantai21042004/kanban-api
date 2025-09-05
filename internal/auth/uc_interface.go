package auth

import (
	"context"

	"github.com/nguyentantai21042004/kanban-api/internal/models"
)

//go:generate mockery --name UseCase
type UseCase interface {
	Login(ctx context.Context, sc models.Scope, ip LoginInput) (LoginOutput, error)
	RefreshToken(ctx context.Context, sc models.Scope, ip RefreshTokenInput) (RefreshTokenOutput, error)
	Logout(ctx context.Context, sc models.Scope) error
}
