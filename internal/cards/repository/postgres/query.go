package postgres

import (
	"context"
	"strings"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"gitlab.com/tantai-kanban/kanban-api/internal/cards"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

func (r implRepository) buildGetQuery(ctx context.Context, fils cards.Filter) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if len(fils.IDs) > 0 {
		for _, id := range fils.IDs {
			if err := postgres.IsUUID(id); err != nil {
				r.l.Errorf(ctx, "internal.labels.repository.postgres.buildGetQuery.InvalidID: %v", err)
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
			r.l.Errorf(ctx, "internal.cards.repository.postgres.buildGetQuery.InvalidBoardID: %v", err)
			return nil, err
		}
		qr = append(qr, qm.Where("board_id = ?", fils.BoardID))
	}

	if fils.ListID != "" {
		if err := postgres.IsUUID(fils.ListID); err != nil {
			r.l.Errorf(ctx, "internal.cards.repository.postgres.buildGetQuery.InvalidListID: %v", err)
			return nil, err
		}
		qr = append(qr, qm.Where("list_id = ?", fils.ListID))
	}

	if fils.Keyword != "" {
		qr = append(qr, qm.Where("Name ILIKE ? OR description ILIKE ?", "%"+fils.Keyword+"%", "%"+fils.Keyword+"%"))
	}

	if fils.AssignedTo != "" {
		if err := postgres.IsUUID(fils.AssignedTo); err != nil {
			r.l.Errorf(ctx, "internal.cards.repository.postgres.buildGetQuery.InvalidAssignedTo: %v", err)
			return nil, err
		}
		qr = append(qr, qm.Where("assigned_to = ?", fils.AssignedTo))
	}

	if fils.Priority != "" {
		qr = append(qr, qm.Where("priority = ?", fils.Priority))
	}

	if len(fils.Tags) > 0 {
		for _, tag := range fils.Tags {
			qr = append(qr, qm.Where("? = ANY(tags)", tag))
		}
	}

	if fils.DueDateFrom != nil {
		qr = append(qr, qm.Where("due_date >= ?", fils.DueDateFrom))
	}

	if fils.DueDateTo != nil {
		qr = append(qr, qm.Where("due_date <= ?", fils.DueDateTo))
	}

	if fils.StartDateFrom != nil {
		qr = append(qr, qm.Where("start_date >= ?", fils.StartDateFrom))
	}

	if fils.StartDateTo != nil {
		qr = append(qr, qm.Where("start_date <= ?", fils.StartDateTo))
	}

	if fils.CompletionDateFrom != nil {
		qr = append(qr, qm.Where("completion_date >= ?", fils.CompletionDateFrom))
	}

	if fils.CompletionDateTo != nil {
		qr = append(qr, qm.Where("completion_date <= ?", fils.CompletionDateTo))
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
