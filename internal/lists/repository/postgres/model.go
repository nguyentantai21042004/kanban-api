package postgres

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists/repository"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

func (r implRepository) buildModel(ctx context.Context, opts repository.CreateOptions) dbmodels.List {
	m := dbmodels.List{
		BoardID:  opts.BoardID,
		Title:    opts.Title,
		Position: types.Decimal{Big: decimal.New(int64(opts.Position), 0)},
	}

	return m
}

func (r implRepository) buildUpdateModel(ctx context.Context, opts repository.UpdateOptions) (dbmodels.List, []string, error) {
	list := dbmodels.List{
		Title:    opts.Title,
		Position: types.Decimal{Big: decimal.New(int64(opts.Position), 0)},
	}
	cols := make([]string, 0)
	cols = append(cols, dbmodels.ListColumns.Title)
	cols = append(cols, dbmodels.ListColumns.Position)

	if err := postgres.IsUUID(opts.ID); err != nil {
		r.l.Errorf(ctx, "internal.lists.repository.postgres.buildUpdateModel.IsUUID: %v", err)
		return dbmodels.List{}, nil, err
	}
	list.ID = opts.ID

	list.UpdatedAt = r.clock()
	cols = append(cols, dbmodels.ListColumns.UpdatedAt)

	return list, cols, nil
}
