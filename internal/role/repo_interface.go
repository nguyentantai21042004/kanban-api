package role

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

//go:generate mockery --name Repository
type Repository interface {
	GetOne(ctx context.Context, sc models.Scope, opts GetOneOptions) (models.Role, error)
	Detail(ctx context.Context, sc models.Scope, ID string) (models.Role, error)
	Get(ctx context.Context, sc models.Scope, opts GetOptions) ([]models.Role, paginator.Paginator, error)
	List(ctx context.Context, sc models.Scope, opts ListOptions) ([]models.Role, error)
}
