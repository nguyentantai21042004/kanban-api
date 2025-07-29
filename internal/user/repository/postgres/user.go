package postgre

import (
	"context"
	"database/sql"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
)

func (r *repository) Detail(ctx context.Context, sc models.Scope, ID string) (models.User, error) {
	// TODO: Implement actual database query
	// For now, return empty user
	return models.User{}, sql.ErrNoRows
}

func (r *repository) Create(ctx context.Context, sc models.Scope, opts user.CreateOptions) (models.User, error) {
	// TODO: Implement actual database query
	// For now, return the user model as is
	return opts.User, nil
}

func (r *repository) Update(ctx context.Context, sc models.Scope, opts user.UpdateOptions) (models.User, error) {
	// TODO: Implement actual database query
	// For now, return the user model as is
	return opts.User, nil
}

func (r *repository) GetByEmail(ctx context.Context, sc models.Scope, email string) (models.User, error) {
	// TODO: Implement actual database query
	// For now, return empty user
	return models.User{}, sql.ErrNoRows
}
