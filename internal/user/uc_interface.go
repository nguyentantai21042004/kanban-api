package user

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

//go:generate mockery --name UseCase
type UseCase interface {
	Detail(ctx context.Context, sc models.Scope, ID string) (UserOutput, error)
	DetailMe(ctx context.Context, sc models.Scope) (UserOutput, error)
	UpdateProfile(ctx context.Context, sc models.Scope, ip UpdateProfileInput) (UserOutput, error)
	Create(ctx context.Context, sc models.Scope, ip CreateInput) (UserOutput, error) // Chỉ Super Admin
	GetOne(ctx context.Context, sc models.Scope, ip GetOneInput) (models.User, error)
}
