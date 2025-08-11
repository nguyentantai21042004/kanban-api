package repository

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

//go:generate mockery --name Repository
type Repository interface {
	GetPosition(ctx context.Context, sc models.Scope, opts GetPositionOptions) (string, error)
	List(ctx context.Context, sc models.Scope, opts ListOptions) ([]models.List, error)
	Get(ctx context.Context, sc models.Scope, opts GetOptions) ([]models.List, paginator.Paginator, error)
	Create(ctx context.Context, sc models.Scope, opts CreateOptions) (models.List, error)
	Update(ctx context.Context, sc models.Scope, opts UpdateOptions) (models.List, error)
	Detail(ctx context.Context, sc models.Scope, id string) (models.List, error)
	Delete(ctx context.Context, sc models.Scope, ids []string) error
	Move(ctx context.Context, sc models.Scope, opts MoveOptions) (models.List, error)
}
