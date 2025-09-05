package postgres

import (
	"context"

	"github.com/aarondl/null/v8"
	"github.com/nguyentantai21042004/kanban-api/internal/comments/repository"
	"github.com/nguyentantai21042004/kanban-api/internal/dbmodels"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/pkg/postgres"
)

func (r implRepository) buildModel(ctx context.Context, sc models.Scope, opts repository.CreateOptions) dbmodels.Comment {
	m := dbmodels.Comment{
		CardID:    opts.CardID,
		UserID:    sc.UserID,
		Content:   opts.Content,
		CreatedAt: r.clock(),
		UpdatedAt: r.clock(),
	}

	if opts.ParentID != nil {
		m.ParentID = null.StringFrom(*opts.ParentID)
	}

	return m
}

func (r implRepository) buildUpdateModel(ctx context.Context, opts repository.UpdateOptions) (dbmodels.Comment, []string, error) {
	comment := dbmodels.Comment{
		Content: opts.Content,
	}
	cols := make([]string, 0)
	cols = append(cols, dbmodels.CommentColumns.Content)

	if err := postgres.IsUUID(opts.ID); err != nil {
		r.l.Errorf(ctx, "internal.comments.repository.postgres.buildUpdateModel.IsUUID: %v", err)
		return dbmodels.Comment{}, nil, err
	}
	comment.ID = opts.ID

	comment.UpdatedAt = r.clock()
	cols = append(cols, dbmodels.CommentColumns.UpdatedAt)

	return comment, cols, nil
}
