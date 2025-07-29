package postgres

import (
	"context"
	"database/sql"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/role"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

func (r *repository) Detail(ctx context.Context, sc models.Scope, ID string) (models.Role, error) {
	// TODO: Implement actual database query
	// For now, return empty role
	return models.Role{}, sql.ErrNoRows
}

func (r *repository) GetOne(ctx context.Context, sc models.Scope, opts role.GetOneOptions) (models.Role, error) {
	// TODO: Implement actual database query
	// For now, return empty role
	return models.Role{}, sql.ErrNoRows
}

func (r *repository) Get(ctx context.Context, sc models.Scope, opts role.GetOptions) ([]models.Role, paginator.Paginator, error) {
	// TODO: Implement actual database query
	// For now, return empty results
	return []models.Role{}, paginator.Paginator{}, nil
}

func (r *repository) List(ctx context.Context, sc models.Scope, opts role.ListOptions) ([]models.Role, error) {
	// TODO: Implement actual database query
	// For now, return empty results
	return []models.Role{}, nil
}
