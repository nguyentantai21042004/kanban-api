package postgres

import (
	"context"

	"github.com/aarondl/null/v8"
	"github.com/nguyentantai21042004/kanban-api/internal/boards/repository"
	"github.com/nguyentantai21042004/kanban-api/internal/dbmodels"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/pkg/postgres"
)

func (r implRepository) buildModel(ctx context.Context, sc models.Scope, opts repository.CreateOptions) dbmodels.Board {
	m := dbmodels.Board{
		Name:        opts.Name,
		Alias:       null.StringFrom(opts.Alias),
		Description: null.StringFrom(opts.Description),
		CreatedBy:   null.StringFrom(sc.UserID),
		CreatedAt:   r.clock(),
		UpdatedAt:   r.clock(),
	}

	return m
}

func (r implRepository) buildUpdateModel(ctx context.Context, opts repository.UpdateOptions) (dbmodels.Board, []string, error) {
	board := dbmodels.Board{
		Name:        opts.Name,
		Alias:       null.StringFrom(opts.Alias),
		Description: null.StringFrom(opts.Description),
	}
	cols := make([]string, 0)
	cols = append(cols, dbmodels.BoardColumns.Name)
	cols = append(cols, dbmodels.BoardColumns.Alias)
	cols = append(cols, dbmodels.BoardColumns.Description)

	if err := postgres.IsUUID(opts.ID); err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.buildUpdateModel.IsUUID: %v", err)
		return dbmodels.Board{}, nil, err
	}
	board.ID = opts.ID

	board.UpdatedAt = r.clock()
	cols = append(cols, dbmodels.BoardColumns.UpdatedAt)

	return board, cols, nil
}
