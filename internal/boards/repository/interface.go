package repository

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

//go:generate mockery --name Repository
type Repository interface {
	Create(ctx context.Context, sc models.Scope, opts CreateOptions) (models.Board, error)
	Detail(ctx context.Context, sc models.Scope, id string) (models.Board, error)
	Update(ctx context.Context, sc models.Scope, opts UpdateOptions) (models.Board, error)
}
