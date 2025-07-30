package postgres

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"gitlab.com/tantai-kanban/kanban-api/internal/user/repository"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

func (r implRepository) buildDetailQuery(ctx context.Context, ID string) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if err := postgres.IsUUID(ID); err != nil {
		r.l.Errorf(ctx, "internal.user.repository.postgres.buildDetailQuery.InvalidID: %v", err)
		return nil, err
	}
	qr = append(qr, qm.Where("id = ?", ID))

	return qr, nil
}

func (r implRepository) buildGetOneQuery(ctx context.Context, opts repository.GetOneOptions) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if opts.Username != "" {
		qr = append(qr, qm.Where("username ILIKE ?", "%"+opts.Username+"%"))
	}

	return qr, nil
}
