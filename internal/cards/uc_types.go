package cards

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type Filter struct {
	IDs     []string
	ListID  string
	Keyword string
}

type GetInput struct {
	Filter   Filter
	PagQuery paginator.PaginateQuery
}

type CreateInput struct {
	ListID      string
	Title       string
	Description string
	Priority    models.CardPriority
	Labels      []string
	DueDate     *time.Time
}

type UpdateInput struct {
	ID          string
	Title       *string
	Description *string
	Priority    *models.CardPriority
	Labels      *[]string
	DueDate     **time.Time
}

type MoveInput struct {
	ID       string
	ListID   string
	Position float64
}

type GetOutput struct {
	Cards      []models.Card
	Pagination paginator.Paginator
}

type DetailOutput struct {
	Card models.Card
}

type GetActivitiesInput struct {
	CardID string
}

type GetActivitiesOutput struct {
	Activities []models.CardActivity
}
