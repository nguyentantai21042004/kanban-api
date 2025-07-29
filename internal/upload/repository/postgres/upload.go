package postgres

import (
	"context"
	"database/sql"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/upload"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

func (r *repository) Create(ctx context.Context, sc models.Scope, opts upload.CreateOptions) (models.Upload, error) {
	// TODO: Implement actual database query
	// For now, return the upload model as is
	return opts.Upload, nil
}

func (r *repository) Detail(ctx context.Context, sc models.Scope, ID string) (models.Upload, error) {
	// TODO: Implement actual database query
	// For now, return empty upload
	return models.Upload{}, sql.ErrNoRows
}

func (r *repository) Get(ctx context.Context, sc models.Scope, opts upload.GetOptions) ([]models.Upload, paginator.Paginator, error) {
	// TODO: Implement actual database query
	// For now, return empty results
	return []models.Upload{}, paginator.Paginator{}, nil
}
