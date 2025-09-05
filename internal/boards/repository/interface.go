package repository

import (
	"context"

	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/pkg/paginator"
)

//go:generate mockery --name Repository
type Repository interface {
	List(ctx context.Context, sc models.Scope, opts ListOptions) ([]models.Board, paginator.Paginator, error)
	Get(ctx context.Context, sc models.Scope, opts GetOptions) ([]models.Board, paginator.Paginator, error)
	Create(ctx context.Context, sc models.Scope, opts CreateOptions) (models.Board, error)
	Update(ctx context.Context, sc models.Scope, opts UpdateOptions) (models.Board, error)
	Detail(ctx context.Context, sc models.Scope, id string) (models.Board, error)
	Delete(ctx context.Context, sc models.Scope, ids []string) error
}
