package postgres

import (
	"context"
	"database/sql"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
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
	// Check database connection first
	if err := r.database.PingContext(ctx); err != nil {
		r.l.Errorf(ctx, "Database connection failed: %v", err)
		return models.User{}, err
	}

	// Convert models.User to dbmodels.User with proper null types
	dbUser := &dbmodels.User{
		ID:           opts.User.ID,
		Username:     opts.User.Username,
		PasswordHash: null.StringFrom(opts.User.PasswordHash),
		FullName:     null.StringFrom(opts.User.FullName),
		RoleID:       null.StringFrom(opts.User.RoleID),
		IsActive:     null.BoolFrom(opts.User.IsActive),
		CreatedAt:    null.TimeFrom(opts.User.CreatedAt),
		UpdatedAt:    null.TimeFrom(opts.User.UpdatedAt),
	}

	r.l.Infof(ctx, "Creating user with ID: %s, Username: %s", dbUser.ID, dbUser.Username)

	// Insert into database
	err := dbUser.Insert(ctx, r.database, boil.Infer())
	if err != nil {
		r.l.Errorf(ctx, "internal.user.repository.postgres.Create.Insert: %v", err)
		return models.User{}, err
	}

	r.l.Infof(ctx, "Successfully created user with ID: %s", dbUser.ID)

	// Verify the user was actually created
	createdUser, err := dbmodels.Users(dbmodels.UserWhere.ID.EQ(dbUser.ID)).One(ctx, r.database)
	if err != nil {
		r.l.Errorf(ctx, "Failed to verify created user: %v", err)
		return models.User{}, err
	}

	r.l.Infof(ctx, "Verified user created in database: %s", createdUser.ID)

	// Return the created user
	return opts.User, nil
}

func (r *implRepository) Update(ctx context.Context, sc models.Scope, opts repository.UpdateOptions) (models.User, error) {
	// Convert models.User to dbmodels.User with proper null types
	dbUser := &dbmodels.User{
		ID:           opts.User.ID,
		Username:     opts.User.Username,
		PasswordHash: null.StringFrom(opts.User.PasswordHash),
		FullName:     null.StringFrom(opts.User.FullName),
		AvatarURL:    null.StringFrom(opts.User.AvatarURL),
		RoleID:       null.StringFrom(opts.User.RoleID),
		IsActive:     null.BoolFrom(opts.User.IsActive),
		CreatedAt:    null.TimeFrom(opts.User.CreatedAt),
		UpdatedAt:    null.TimeFrom(opts.User.UpdatedAt),
	}

	// Update in database
	rowsAffected, err := dbUser.Update(ctx, r.database, boil.Infer())
	if err != nil {
		r.l.Errorf(ctx, "internal.user.repository.postgres.Update.Update: %v", err)
		return models.User{}, err
	}

	if rowsAffected == 0 {
		r.l.Warnf(ctx, "internal.user.repository.postgres.Update.Update: no rows affected")
		return models.User{}, repository.ErrNotFound
	}

	// Return the updated user
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
