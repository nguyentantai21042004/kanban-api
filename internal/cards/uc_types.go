package cards

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type Filter struct {
	IDs                []string
	ListID             string
	BoardID            string
	Keyword            string
	AssignedTo         string
	Priority           models.CardPriority
	Tags               []string
	DueDateFrom        *time.Time
	DueDateTo          *time.Time
	StartDateFrom      *time.Time
	StartDateTo        *time.Time
	CompletionDateFrom *time.Time
	CompletionDateTo   *time.Time
}

type GetInput struct {
	Filter   Filter
	PagQuery paginator.PaginateQuery
}

type ChecklistItemInput struct {
	Content     string `json:"content"`
	IsCompleted bool   `json:"is_completed"`
}

type CreateInput struct {
	BoardID        string
	ListID         string
	Name           string
	Description    string
	Priority       models.CardPriority
	Labels         []string
	DueDate        *time.Time
	AssignedTo     *string
	EstimatedHours *float64
	StartDate      *time.Time
	Tags           []string
	Checklist      []ChecklistItemInput
}

type UpdateInput struct {
	ID             string
	Name           string
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
}

type MoveInput struct {
	AfterID  string
	ID       string
	BeforeID string
	ListID   string
}

type GetOutput struct {
	Cards      []models.Card
	Pagination paginator.Paginator
}

type DetailOutput struct {
	Card  models.Card
	List  models.List
	Board models.Board
	Users []models.User
}

type GetActivitiesInput struct {
	CardID   string
	PagQuery paginator.PaginateQuery
}

type GetActivitiesOutput struct {
	Card       models.Card
	Activities []models.CardActivity
	Pagination paginator.Paginator
}

type AssignInput struct {
	CardID     string
	AssignedTo string
}

type UnassignInput struct {
	CardID string
}

type AddAttachmentInput struct {
	CardID       string
	AttachmentID string
}

type RemoveAttachmentInput struct {
	CardID       string
	AttachmentID string
}

type UpdateTimeTrackingInput struct {
	CardID         string
	EstimatedHours *float64
	ActualHours    *float64
}

type UpdateChecklistInput struct {
	CardID    string
	Checklist []models.ChecklistItem
}

type AddTagInput struct {
	CardID string
	Tag    string
}

type RemoveTagInput struct {
	CardID string
	Tag    string
}

type SetStartDateInput struct {
	CardID    string
	StartDate *time.Time
}

type SetCompletionDateInput struct {
	CardID         string
	CompletionDate *time.Time
}
