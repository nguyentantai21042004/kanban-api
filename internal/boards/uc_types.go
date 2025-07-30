package boards

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type Filter struct {
	IDs       []string
	Keyword   string
	CreatedBy string
}

type GetInput struct {
	Filter   Filter
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
	Users      []models.User
	Pagination paginator.Paginator
}

type DetailOutput struct {
	Board models.Board
	User  models.User
}
