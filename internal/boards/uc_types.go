package boards

import (
	"time"

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
	Users []models.User
}

type BoardWithDetailsOutput struct {
	Board  models.Board
	Lists  []models.List
	Cards  []models.Card
	Labels []models.Label
}

// Dashboard aggregation for boards
type DashboardInput struct {
	From time.Time
	To   time.Time
}

type BoardsDashboardOutput struct {
	Total  int64
	Active int64
}
