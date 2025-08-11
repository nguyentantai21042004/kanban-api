package lists

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
	BoardID  string
	Name     string
	Position string
}

type UpdateInput struct {
	ID       string
	Name     string
	Position string
}

type GetOutput struct {
	Lists      []models.List
	Pagination paginator.Paginator
}

type DetailOutput struct {
	List models.List
}
