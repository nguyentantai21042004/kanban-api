package boards

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/boards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type GetInput struct {
	Filter   repository.Filter
	PagQuery paginator.PaginateQuery
}

type CreateInput struct {
	Name        string
	Description string
}

type UpdateInput struct {
	ID          string
	Name        string
	Description string
}

type GetOutput struct {
	Boards     []models.Board
	Pagination paginator.Paginator
}

type DetailOutput struct {
	Board models.Board
}
