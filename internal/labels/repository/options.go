package repository

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/labels"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type GetOptions struct {
	Filter   labels.Filter
	PagQuery paginator.PaginateQuery
}

type CreateOptions struct {
	BoardID string
	Name    string
	Color   string
}

type UpdateOptions struct {
	ID       string
	Name     string
	Color    string
	OldModel models.Label
}
