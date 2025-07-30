package postgres

import (
	"context"
	"strings"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"gitlab.com/tantai-kanban/kanban-api/internal/boards"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
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
		placeholders := make([]string, len(fils.IDs))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		qr = append(qr, qm.WhereIn("id IN ("+strings.Join(placeholders, ",")+")", postgres.ConvertToInterface(fils.IDs)...))
	}

	if fils.Keyword != "" {
		qr = append(qr, qm.Where("name ILIKE ?", "%"+fils.Keyword+"%"))
		qr = append(qr, qm.Or("alias ILIKE ?", "%"+fils.Keyword+"%"))
		qr = append(qr, qm.Or("description ILIKE ?", "%"+fils.Keyword+"%"))
		qr = append(qr, qm.Or("id IN (SELECT board_id FROM lists WHERE title ILIKE ?)", "%"+fils.Keyword+"%"))
		qr = append(qr, qm.Or("id IN (SELECT board_id FROM cards WHERE title ILIKE ?)", "%"+fils.Keyword+"%"))
		qr = append(qr, qm.Or("id IN (SELECT board_id FROM cards WHERE description ILIKE ?)", "%"+fils.Keyword+"%"))
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
