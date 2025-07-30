package repository

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

//go:generate mockery --name Repository
type Repository interface {
	Detail(ctx context.Context, sc models.Scope, ID string) (models.User, error)
	Create(ctx context.Context, sc models.Scope, opts CreateOptions) (models.User, error)
	Update(ctx context.Context, sc models.Scope, opts UpdateOptions) (models.User, error)
	GetOne(ctx context.Context, sc models.Scope, ip GetOneOptions) (models.User, error)
}
