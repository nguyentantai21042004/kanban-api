package postgres

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/role/repository"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

func (r implRepository) buildDetailQuery(ctx context.Context, ID string) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if err := postgres.IsUUID(ID); err != nil {
		r.l.Errorf(ctx, "internal.role.repository.postgres.buildDetailQuery.InvalidID: %v", err)
		return nil, err
	}
	qr = append(qr, qm.Where("id = ?", ID))

	return qr, nil
}

func (r implRepository) buildListQuery(ctx context.Context, sc models.Scope, opts repository.ListOptions) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if len(opts.Filter.IDs) > 0 {
		for _, id := range opts.Filter.IDs {
			if err := postgres.IsUUID(id); err != nil {
				r.l.Errorf(ctx, "internal.role.repository.postgres.buildListQuery.InvalidID: %v", err)
				return nil, err
			}
		}
		qr = append(qr, dbmodels.RoleWhere.ID.IN(opts.Filter.IDs))
	}

	if opts.Filter.Code != "" {
		qr = append(qr, qm.Where("code ~* ?", opts.Filter.Code))
	}

	if opts.Filter.IsActive {
		qr = append(qr, qm.Where("is_active = ?", opts.Filter.IsActive))
	}

	return qr, nil
}
