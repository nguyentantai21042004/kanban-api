package postgres

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"gitlab.com/tantai-kanban/kanban-api/internal/user/repository"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

func (r implRepository) buildGetOneQuery(ctx context.Context, opts repository.GetOneOptions) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if opts.Username != "" {
		qr = append(qr, qm.Where("username ILIKE ?", "%"+opts.Username+"%"))
	}

	return qr, nil
}
