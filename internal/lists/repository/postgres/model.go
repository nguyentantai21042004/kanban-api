package postgres

import (
	"context"

	"github.com/aarondl/null/v8"
	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

func (r implRepository) buildModel(ctx context.Context, sc models.Scope, opts repository.CreateOptions) dbmodels.List {
	m := dbmodels.List{
		BoardID:   opts.BoardID,
		Name:      opts.Name,
		Position:  opts.Position,
		CreatedBy: null.StringFrom(sc.UserID),
	}

	return m
}

func (r implRepository) buildUpdateModel(ctx context.Context, opts repository.UpdateOptions) (dbmodels.List, []string, error) {
	list := dbmodels.List{
		Name: opts.Name,
	}
	cols := make([]string, 0)
	cols = append(cols, dbmodels.ListColumns.Name)

	if err := postgres.IsUUID(opts.ID); err != nil {
		r.l.Errorf(ctx, "internal.lists.repository.postgres.buildUpdateModel.IsUUID: %v", err)
		return dbmodels.List{}, nil, err
	}
	list.ID = opts.ID

	list.UpdatedAt = r.clock()
	cols = append(cols, dbmodels.ListColumns.UpdatedAt)

	return list, cols, nil
}
