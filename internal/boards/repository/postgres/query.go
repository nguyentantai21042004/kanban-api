package postgres

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/nguyentantai21042004/kanban-api/internal/boards"
	"github.com/nguyentantai21042004/kanban-api/internal/dbmodels"
	"github.com/nguyentantai21042004/kanban-api/pkg/postgres"
)

func (r implRepository) buildGetQuery(ctx context.Context, fils boards.Filter) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if len(fils.IDs) > 0 {
		for _, id := range fils.IDs {
			if err := postgres.IsUUID(id); err != nil {
				r.l.Errorf(ctx, "internal.boards.repository.postgres.buildGetQuery.InvalidID: %v", err)
				return nil, err
			}
		}
		qr = append(qr, dbmodels.BoardWhere.ID.IN(fils.IDs))
	}

	if fils.Keyword != "" {
		qr = append(qr, qm.Where("alias ILIKE ? OR description ILIKE ?", "%"+fils.Keyword+"%", "%"+fils.Keyword+"%"))
	}

	if fils.CreatedBy != "" {
		if err := postgres.IsUUID(fils.CreatedBy); err != nil {
			r.l.Errorf(ctx, "internal.boards.repository.postgres.buildGetQuery.InvalidCreatedBy: %v", err)
			return nil, err
		}
		qr = append(qr, qm.Where("created_by = ?", fils.CreatedBy))
	}

	return qr, nil
}

func (r implRepository) buildDetailQuery(ctx context.Context, ID string) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if err := postgres.IsUUID(ID); err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.buildDetailQuery.InvalidID: %v", err)
		return nil, err
	}
	qr = append(qr, qm.Where("id = ?", ID))

	return qr, nil
}

func (r implRepository) buildDeleteQuery(ctx context.Context, IDs []string) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	for _, ID := range IDs {
		if err := postgres.IsUUID(ID); err != nil {
			r.l.Errorf(ctx, "internal.boards.repository.postgres.buildDeleteQuery.InvalidID: %v", err)
			return nil, err
		}
		qr = append(qr, qm.Where("id = ?", ID))
	}

	return qr, nil
}
