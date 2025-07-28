package postgres

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/internal/labels/repository"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

func (r implRepository) buildModel(ctx context.Context, opts repository.CreateOptions) dbmodels.Label {
	m := dbmodels.Label{
		BoardID: opts.BoardID,
		Name:    opts.Name,
		Color:   opts.Color,
	}

	return m
}

func (r implRepository) buildUpdateModel(ctx context.Context, opts repository.UpdateOptions) (dbmodels.Label, []string, error) {
	label := dbmodels.Label{
		Name:  opts.Name,
		Color: opts.Color,
	}
	cols := make([]string, 0)
	cols = append(cols, dbmodels.LabelColumns.Name)
	cols = append(cols, dbmodels.LabelColumns.Color)

	if err := postgres.IsUUID(opts.ID); err != nil {
		r.l.Errorf(ctx, "internal.labels.repository.postgres.buildUpdateModel.IsUUID: %v", err)
		return dbmodels.Label{}, nil, err
	}
	label.ID = opts.ID

	return label, cols, nil
}
