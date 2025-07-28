package labels

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type Filter struct {
	IDs     []string
	BoardID string
	Keyword string
}

type GetInput struct {
	Filter   Filter
	PagQuery paginator.PaginateQuery
}

type CreateInput struct {
	BoardID string
	Name    string
	Color   string
}

type UpdateInput struct {
	ID     string
	Name   string
	Color  string
}

type GetOutput struct {
	Labels     []models.Label
	Pagination paginator.Paginator
}

type DetailOutput struct {
	Label models.Label
}
