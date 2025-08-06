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
	ListID         string
	Name           string
	Description    string
	Position       float64
	Priority       models.CardPriority
	Labels         []string
	DueDate        *time.Time
	CreatedBy      string
	AssignedTo     *string
	EstimatedHours *float64
	StartDate      *time.Time
	Tags           []string
	Checklist      []models.ChecklistItem
}

type UpdateOptions struct {
	ID             string
	Name           *string
	Description    *string
	Priority       *models.CardPriority
	Labels         *[]string
	DueDate        *time.Time
	AssignedTo     *string
	EstimatedHours *float64
	ActualHours    *float64
	StartDate      *time.Time
	CompletionDate *time.Time
	Tags           *[]string
	Checklist      *[]models.ChecklistItem
	OldModel       models.Card
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

// New option types for enhanced functionality
type AssignOptions struct {
	CardID     string
	AssignedTo string
	OldModel   models.Card
}

type UnassignOptions struct {
	CardID   string
	OldModel models.Card
}

type AddAttachmentOptions struct {
	CardID       string
	AttachmentID string
	OldModel     models.Card
}

type RemoveAttachmentOptions struct {
	CardID       string
	AttachmentID string
	OldModel     models.Card
}

type UpdateTimeTrackingOptions struct {
	CardID         string
	EstimatedHours *float64
	ActualHours    *float64
	OldModel       models.Card
}

type UpdateChecklistOptions struct {
	CardID    string
	Checklist []models.ChecklistItem
	OldModel  models.Card
}

type AddTagOptions struct {
	CardID   string
	Tag      string
	OldModel models.Card
}

type RemoveTagOptions struct {
	CardID   string
	Tag      string
	OldModel models.Card
}

type SetStartDateOptions struct {
	CardID    string
	StartDate *time.Time
	OldModel  models.Card
}

type SetCompletionDateOptions struct {
	CardID         string
	CompletionDate *time.Time
	OldModel       models.Card
}
