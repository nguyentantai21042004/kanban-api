package postgres

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/nguyentantai21042004/kanban-api/internal/comments"
	"github.com/nguyentantai21042004/kanban-api/internal/dbmodels"
	"github.com/nguyentantai21042004/kanban-api/pkg/postgres"
)

func (r implRepository) buildGetQuery(ctx context.Context, fils comments.Filter) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if len(fils.IDs) > 0 {
		for _, id := range fils.IDs {
			if err := postgres.IsUUID(id); err != nil {
				r.l.Errorf(ctx, "internal.comments.repository.postgres.buildGetQuery.InvalidID: %v", err)
				return nil, err
			}
		}
		qr = append(qr, dbmodels.CommentWhere.ID.IN(fils.IDs))
	}

	if fils.Keyword != "" {
		qr = append(qr, qm.Where("content ILIKE ?", "%"+fils.Keyword+"%"))
	}

	if fils.CardID != "" {
		if err := postgres.IsUUID(fils.CardID); err != nil {
			r.l.Errorf(ctx, "internal.comments.repository.postgres.buildGetQuery.InvalidCardID: %v", err)
			return nil, err
		}
		qr = append(qr, qm.Where("card_id = ?", fils.CardID))
	}

	if fils.UserID != "" {
		if err := postgres.IsUUID(fils.UserID); err != nil {
			r.l.Errorf(ctx, "internal.comments.repository.postgres.buildGetQuery.InvalidUserID: %v", err)
			return nil, err
		}
		qr = append(qr, qm.Where("user_id = ?", fils.UserID))
	}

	if fils.ParentID != "" {
		if err := postgres.IsUUID(fils.ParentID); err != nil {
			r.l.Errorf(ctx, "internal.comments.repository.postgres.buildGetQuery.InvalidParentID: %v", err)
			return nil, err
		}
		qr = append(qr, qm.Where("parent_id = ?", fils.ParentID))
	}

	return qr, nil
}

func (r implRepository) buildDetailQuery(ctx context.Context, ID string) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if err := postgres.IsUUID(ID); err != nil {
		r.l.Errorf(ctx, "internal.comments.repository.postgres.buildDetailQuery.InvalidID: %v", err)
		return nil, err
	}
	qr = append(qr, qm.Where("id = ?", ID))

	return qr, nil
}

func (r implRepository) buildDeleteQuery(ctx context.Context, IDs []string) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	for _, ID := range IDs {
		if err := postgres.IsUUID(ID); err != nil {
			r.l.Errorf(ctx, "internal.comments.repository.postgres.buildDeleteQuery.InvalidID: %v", err)
			return nil, err
		}
		qr = append(qr, qm.Where("id = ?", ID))
	}

	return qr, nil
}

func (r implRepository) buildGetByCardQuery(ctx context.Context, cardID string) ([]qm.QueryMod, error) {
	qr := postgres.BuildQueryWithSoftDelete()

	if err := postgres.IsUUID(cardID); err != nil {
		r.l.Errorf(ctx, "internal.comments.repository.postgres.buildGetByCardQuery.InvalidCardID: %v", err)
		return nil, err
	}
	qr = append(qr, qm.Where("card_id = ?", cardID))
	qr = append(qr, qm.OrderBy("created_at ASC"))

	return qr, nil
}
