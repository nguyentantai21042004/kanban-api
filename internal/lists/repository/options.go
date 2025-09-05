package repository

import (
	"github.com/nguyentantai21042004/kanban-api/internal/lists"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/pkg/paginator"
)

type GetPositionOptions struct {
	BoardID string
	ASC     bool
}

type ListOptions struct {
	Filter lists.Filter
}

type GetOptions struct {
	Filter   lists.Filter
	PagQuery paginator.PaginateQuery
}

type CreateOptions struct {
	BoardID  string
	Name     string
	Position string
}

type UpdateOptions struct {
	ID       string
	Name     string
	OldModel models.List
}

type MoveOptions struct {
	ID          string
	BoardID     string
	NewPosition string
}
