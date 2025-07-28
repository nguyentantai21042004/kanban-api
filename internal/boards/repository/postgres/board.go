package postgres

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/boil"
	"gitlab.com/tantai-kanban/kanban-api/internal/boards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

func (r implRepository) Create(ctx context.Context, sc models.Scope, opts repository.CreateOptions) (models.Board, error) {
	m := r.buildModel(ctx, opts)

	err := m.Insert(ctx, r.database, boil.Infer())
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Create.Insert: %v", err)
		return models.Board{}, err
	}

	return models.NewBoard(m), nil
}

func (r implRepository) Update(ctx context.Context, sc models.Scope, opts repository.UpdateOptions) (models.Board, error) {
	b, col, err := r.buildUpdateModel(ctx, opts)
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Update.buildUpdateModel: %v", err)
		return models.Board{}, err
	}

	_, err = b.Update(ctx, r.database, boil.Whitelist(col...))
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Update.Update: %v", err)
		return models.Board{}, err
	}

	return models.NewBoard(b), nil
}

func (r implRepository) Detail(ctx context.Context, sc models.Scope, ID string) (models.Board, error) {
	qr, err := r.buildDetailQuery(ctx, ID)
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Detail.buildDetailQuery: %v", err)
		return models.Board{}, err
	}

	board, err := dbmodels.Boards(qr...).One(ctx, r.database)
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Detail.One: %v", err)
		return models.Board{}, err
	}

	return models.NewBoard(*board), nil
}

func (r implRepository) Delete(ctx context.Context, sc models.Scope, IDs []string) error {
	qr, err := r.buildDeleteQuery(ctx, IDs)
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Delete.buildDeleteQuery: %v", err)
		return err
	}

	_, err = dbmodels.Boards(qr...).DeleteAll(ctx, r.database, true)
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Delete.DeleteAll: %v", err)
		return err
	}

	return nil
}
