package repository

import (
	"github.com/nguyentantai21042004/kanban-api/internal/comments"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/pkg/paginator"
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
