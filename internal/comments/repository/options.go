package repository

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/comments"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type GetOptions struct {
	Filter   comments.Filter
	PagQuery paginator.PaginateQuery
}

type CreateOptions struct {
	CardID   string
	Content  string
	ParentID *string
}

type UpdateOptions struct {
	ID       string
	Content  string
	OldModel models.Comment
}
