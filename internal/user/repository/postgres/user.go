package postgres

import (
	"context"
	"database/sql"

	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/user/repository"
)

func (r *implRepository) Detail(ctx context.Context, sc models.Scope, ID string) (models.User, error) {
	qr, err := r.buildDetailQuery(ctx, ID)
	if err != nil {
		r.l.Errorf(ctx, "internal.user.repository.postgres.Detail.buildDetailQuery: %v", err)
		return models.User{}, err
	}

	u, err := dbmodels.Users(qr...).One(ctx, r.database)
	if err != nil {
		if err == sql.ErrNoRows {
			r.l.Errorf(ctx, "internal.user.repository.postgres.Detail.One.NoRows: %v", err)
			return models.User{}, repository.ErrNotFound
		}
		r.l.Errorf(ctx, "internal.user.repository.postgres.Detail.One: %v", err)
		return models.User{}, err
	}

	return *models.NewUser(u), nil
}

func (r *implRepository) List(ctx context.Context, sc models.Scope, opts repository.ListOptions) ([]models.User, error) {
	qr, err := r.buildListQuery(ctx, opts)
	if err != nil {
		r.l.Errorf(ctx, "internal.user.repository.postgres.List.buildListQuery: %v", err)
		return nil, err
	}

	users, err := dbmodels.Users(qr...).All(ctx, r.database)
	if err != nil {
		r.l.Errorf(ctx, "internal.user.repository.postgres.List.All: %v", err)
		return nil, err
	}

	results := make([]models.User, len(users))
	for i, u := range users {
		results[i] = *models.NewUser(u)
	}

	return results, nil
}

func (r *implRepository) Create(ctx context.Context, sc models.Scope, opts repository.CreateOptions) (models.User, error) {
	// TODO: Implement actual database query
	// For now, return the user model as is
	return opts.User, nil
}

func (r *implRepository) Update(ctx context.Context, sc models.Scope, opts repository.UpdateOptions) (models.User, error) {
	// TODO: Implement actual database query
	// For now, return the user model as is
	return opts.User, nil
}

func (r *implRepository) GetOne(ctx context.Context, sc models.Scope, opts repository.GetOneOptions) (models.User, error) {
	qr, err := r.buildGetOneQuery(ctx, opts)
	if err != nil {
		r.l.Errorf(ctx, "internal.user.repository.postgres.GetOne.buildGetOneQuery: %v", err)
		return models.User{}, err
	}

	u, err := dbmodels.Users(qr...).One(ctx, r.database)
	if err != nil {
		if err == sql.ErrNoRows {
			r.l.Errorf(ctx, "internal.user.repository.postgres.GetOne.One.ErrNoRows")
			return models.User{}, repository.ErrNotFound
		}
		r.l.Errorf(ctx, "internal.user.repository.postgres.GetOne.One: %v", err)
		return models.User{}, err
	}

	return *models.NewUser(u), nil
}
