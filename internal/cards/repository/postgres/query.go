package postgres

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"gitlab.com/tantai-kanban/kanban-api/internal/cards"
	"gitlab.com/tantai-kanban/kanban-api/internal/cards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

func (r implRepository) buildGetPositionQuery(opts repository.GetPositionOptions) ([]qm.QueryMod, error) {
	order := "DESC"
	if opts.ASC {
		order = "ASC"
	}

	return []qm.QueryMod{
		dbmodels.CardWhere.ListID.EQ(opts.ListID),
		qm.OrderBy("position " + order),
	}, nil
}

func (r implRepository) buildGetQuery(ctx context.Context, fils cards.Filter) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if len(fils.IDs) > 0 {
		for _, id := range fils.IDs {
			if err := postgres.IsUUID(id); err != nil {
				r.l.Errorf(ctx, "internal.labels.repository.postgres.buildGetQuery.InvalidID: %v", err)
				return nil, err
			}
		}
		qr = append(qr, dbmodels.CardWhere.ID.IN(fils.IDs))
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

	if fils.CreatedBy != "" {
		if err := postgres.IsUUID(fils.CreatedBy); err != nil {
			r.l.Errorf(ctx, "internal.cards.repository.postgres.buildGetQuery.InvalidCreatedBy: %v", err)
			return nil, err
		}
		qr = append(qr, qm.Where("created_by = ?", fils.CreatedBy))
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

	if fils.CreatedFrom != nil {
		qr = append(qr, qm.Where("created_at >= ?", fils.CreatedFrom))
	}

	if fils.CreatedTo != nil {
		qr = append(qr, qm.Where("created_at <= ?", fils.CreatedTo))
	}

	if fils.UpdatedFrom != nil {
		qr = append(qr, qm.Where("updated_at >= ?", fils.UpdatedFrom))
	}

	if fils.UpdatedTo != nil {
		qr = append(qr, qm.Where("updated_at <= ?", fils.UpdatedTo))
	}

	if fils.UncompletedOnly {
		qr = append(qr, qm.Where("completion_date is null"))
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

func (r implRepository) buildGetActivitiesQuery(ctx context.Context, fils repository.ActivityFilter) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	qr = append(qr, qm.Where("card_id = ?", fils.CardID))

	return qr, nil
}
