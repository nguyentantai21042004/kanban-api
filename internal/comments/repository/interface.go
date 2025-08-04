package repository

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

//go:generate mockery --name Repository
type Repository interface {
	Get(ctx context.Context, sc models.Scope, opts GetOptions) ([]models.Comment, paginator.Paginator, error)
	Create(ctx context.Context, sc models.Scope, opts CreateOptions) (models.Comment, error)
	Update(ctx context.Context, sc models.Scope, opts UpdateOptions) (models.Comment, error)
	Detail(ctx context.Context, sc models.Scope, id string) (models.Comment, error)
	Delete(ctx context.Context, sc models.Scope, ids []string) error
	GetByCard(ctx context.Context, sc models.Scope, cardID string) ([]models.Comment, error)
}
