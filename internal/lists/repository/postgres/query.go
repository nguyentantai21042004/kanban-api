package postgres

import (
	"context"
	"strings"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

func (r implRepository) buildGetQuery(ctx context.Context, fils lists.Filter) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if len(fils.IDs) > 0 {
		for _, id := range fils.IDs {
			if err := postgres.IsUUID(id); err != nil {
				r.l.Errorf(ctx, "internal.lists.repository.postgres.buildGetQuery.InvalidID: %v", err)
				return nil, err
			}
		}
		placeholders := make([]string, len(fils.IDs))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		qr = append(qr, qm.WhereIn("id IN ("+strings.Join(placeholders, ",")+")", postgres.ConvertToInterface(fils.IDs)...))
	}

	if fils.BoardID != "" {
		if err := postgres.IsUUID(fils.BoardID); err != nil {
			r.l.Errorf(ctx, "internal.lists.repository.postgres.buildGetQuery.InvalidBoardID: %v", err)
			return nil, err
		}
		qr = append(qr, qm.Where("board_id = ?", fils.BoardID))
	}

	if fils.Keyword != "" {
		qr = append(qr, qm.Where("title ILIKE ?", "%"+fils.Keyword+"%"))
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
