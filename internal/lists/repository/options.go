package repository

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/lists"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type GetOptions struct {
	Filter   lists.Filter
	PagQuery paginator.PaginateQuery
}

type CreateOptions struct {
	BoardID  string
	Name     string
	Position float64
}

type UpdateOptions struct {
	ID       string
	Name     string
	Position float64
	OldModel models.List
}
