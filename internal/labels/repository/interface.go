package repository

import (
	"context"

	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/pkg/paginator"
)

//go:generate mockery --name Repository
type Repository interface {
	Get(ctx context.Context, sc models.Scope, opts GetOptions) ([]models.Label, paginator.Paginator, error)
	Create(ctx context.Context, sc models.Scope, opts CreateOptions) (models.Label, error)
	Update(ctx context.Context, sc models.Scope, opts UpdateOptions) (models.Label, error)
	Detail(ctx context.Context, sc models.Scope, id string) (models.Label, error)
	Delete(ctx context.Context, sc models.Scope, ids []string) error
}
