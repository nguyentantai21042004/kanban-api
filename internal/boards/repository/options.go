package repository

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/boards"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type ListOptions struct {
	Filter boards.Filter
}

type GetOptions struct {
	Filter   boards.Filter
	PagQuery paginator.PaginateQuery
}

type CreateOptions struct {
	Name        string
	Description string
	Alias       string
}

type UpdateOptions struct {
	ID          string
	Name        string
	Description string
	Alias       string
	OldModel    models.Board
}
