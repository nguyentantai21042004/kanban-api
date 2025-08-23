package postgres

import (
	"context"
	"database/sql"

	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/role/repository"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

func (r *implRepository) Detail(ctx context.Context, sc models.Scope, ID string) (models.Role, error) {
	qr, err := r.buildDetailQuery(ctx, ID)
	if err != nil {
		r.l.Errorf(ctx, "internal.role.repository.postgres.Detail.buildDetailQuery: %v", err)
		return models.Role{}, err
	}

	rl, err := dbmodels.Roles(qr...).One(ctx, r.database)
	if err != nil {
		if err == sql.ErrNoRows {
			r.l.Errorf(ctx, "internal.role.repository.postgres.Detail.One.NoRows: %v", err)
			return models.Role{}, repository.ErrNotFound
		}
		r.l.Errorf(ctx, "internal.role.repository.postgres.Detail.One: %v", err)
		return models.Role{}, err
	}

	return models.NewRole(*rl), nil
}

func (r *implRepository) GetOne(ctx context.Context, sc models.Scope, opts repository.GetOneOptions) (models.Role, error) {
	// TODO: Implement actual database query
	// For now, return empty role
	return models.Role{}, sql.ErrNoRows
}

func (r *implRepository) Get(ctx context.Context, sc models.Scope, opts repository.GetOptions) ([]models.Role, paginator.Paginator, error) {
	// TODO: Implement actual database query
	// For now, return empty results
	return []models.Role{}, paginator.Paginator{}, nil
}

func (r *implRepository) List(ctx context.Context, sc models.Scope, opts repository.ListOptions) ([]models.Role, error) {
	qr, err := r.buildListQuery(ctx, sc, opts)
	if err != nil {
		r.l.Errorf(ctx, "internal.role.repository.postgres.List.buildListQuery: %v", err)
		return nil, err
	}

	rl, err := dbmodels.Roles(qr...).All(ctx, r.database)
	if err != nil {
		r.l.Errorf(ctx, "internal.role.repository.postgres.List.All: %v", err)
		return nil, err
	}

	roles := make([]models.Role, len(rl))
	for i, r := range rl {
		roles[i] = models.NewRole(*r)
	}

	return roles, nil
}
