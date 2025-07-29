package repository

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/cards"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type GetOptions struct {
	Filter   cards.Filter
	PagQuery paginator.PaginateQuery
}

type CreateOptions struct {
	ListID      string
	Title       string
	Description string
	Position    float64
	Priority    models.CardPriority
	Labels      []string
	DueDate     *time.Time
}

type UpdateOptions struct {
	ID          string
	Title       *string
	Description *string
	Priority    *models.CardPriority
	Labels      *[]string
	DueDate     **time.Time
	OldModel    models.Card
}

type MoveOptions struct {
	ID       string
	ListID   string
	Position float64
	OldModel models.Card
}

type GetActivitiesOptions struct {
	CardID string
}
