package postgres

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists/repository"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

func (r implRepository) buildGetPositionQuery(opts repository.GetPositionOptions) ([]qm.QueryMod, error) {
	order := "DESC"
	if opts.ASC {
		order = "ASC"
	}

	return []qm.QueryMod{
		dbmodels.ListWhere.BoardID.EQ(opts.BoardID),
		qm.OrderBy("position " + order),
	}, nil
}

func (r implRepository) buildGetQuery(ctx context.Context, fils lists.Filter) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if len(fils.IDs) > 0 {
		for _, id := range fils.IDs {
			if err := postgres.IsUUID(id); err != nil {
				r.l.Errorf(ctx, "internal.lists.repository.postgres.buildGetQuery.InvalidID: %v", err)
				return nil, err
			}
		}
		qr = append(qr, dbmodels.ListWhere.ID.IN(fils.IDs))
	}

	if fils.BoardID != "" {
		if err := postgres.IsUUID(fils.BoardID); err != nil {
			r.l.Errorf(ctx, "internal.lists.repository.postgres.buildGetQuery.InvalidBoardID: %v", err)
			return nil, err
		}
		qr = append(qr, qm.Where("board_id = ?", fils.BoardID))
	}

	if fils.Keyword != "" {
		qr = append(qr, qm.Where("Name ILIKE ?", "%"+fils.Keyword+"%"))
	}

	if fils.CreatedBy != "" {
		if err := postgres.IsUUID(fils.CreatedBy); err != nil {
			r.l.Errorf(ctx, "internal.lists.repository.postgres.buildGetQuery.InvalidCreatedBy: %v", err)
			return nil, err
		}
		qr = append(qr, qm.Where("created_by = ?", fils.CreatedBy))
	}

	return qr, nil
}

func (r implRepository) buildDetailQuery(ctx context.Context, ID string) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if err := postgres.IsUUID(ID); err != nil {
		r.l.Errorf(ctx, "internal.lists.repository.postgres.buildDetailQuery.InvalidID: %v", err)
		return nil, err
	}
	qr = append(qr, qm.Where("id = ?", ID))

	return qr, nil
}

func (r implRepository) buildDeleteQuery(ctx context.Context, IDs []string) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	for _, ID := range IDs {
		if err := postgres.IsUUID(ID); err != nil {
			r.l.Errorf(ctx, "internal.lists.repository.postgres.buildDeleteQuery.InvalidID: %v", err)
			return nil, err
		}
		qr = append(qr, qm.Where("id = ?", ID))
	}

	return qr, nil
}
