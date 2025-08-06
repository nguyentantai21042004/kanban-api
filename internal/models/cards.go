package models

import (
	"encoding/json"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
)

type Card struct {
	ID          string       `json:"id"`
	ListID      string       `json:"list_id"`
	Name        string       `json:"name"`
	Alias       string       `json:"alias"`
	Description string       `json:"description,omitempty"`
	Position    float64      `json:"position"`
	DueDate     *time.Time   `json:"due_date,omitempty"`
	Priority    CardPriority `json:"priority"`
	Labels      []string     `json:"labels,omitempty"`
	IsArchived  bool         `json:"is_archived"`
	CreatedBy   *string      `json:"created_by,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   *time.Time   `json:"deleted_at,omitempty"`
	// New fields from migration
	AssignedTo     *string         `json:"assigned_to,omitempty"`
	Attachments    []string        `json:"attachments,omitempty"`
	EstimatedHours *float64        `json:"estimated_hours,omitempty"`
	ActualHours    *float64        `json:"actual_hours,omitempty"`
	StartDate      *time.Time      `json:"start_date,omitempty"`
	CompletionDate *time.Time      `json:"completion_date,omitempty"`
	Tags           []string        `json:"tags,omitempty"`
	Checklist      []ChecklistItem `json:"checklist,omitempty"`
	LastActivityAt *time.Time      `json:"last_activity_at,omitempty"`
	UpdatedBy      *string         `json:"updated_by,omitempty"`
}

type CardPriority string

const (
	CardPriorityLow    CardPriority = "low"
	CardPriorityMedium CardPriority = "medium"
	CardPriorityHigh   CardPriority = "high"
)

type ChecklistItem struct {
	ID          string    `json:"id"`
	Content     string    `json:"content"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewCard(dbCard dbmodels.Card) Card {
	labels := []string{}
	if dbCard.Labels.Valid {
		_ = json.Unmarshal(dbCard.Labels.JSON, &labels)
	}

	attachments := []string{}
	if dbCard.Attachments.Valid {
		_ = json.Unmarshal(dbCard.Attachments.JSON, &attachments)
	}

	checklist := []ChecklistItem{}
	if dbCard.Checklist.Valid {
		_ = json.Unmarshal(dbCard.Checklist.JSON, &checklist)
	}

	var desc string
	if dbCard.Description.Valid {
		desc = dbCard.Description.String
	}

	var due *time.Time
	if dbCard.DueDate.Valid {
		d := dbCard.DueDate.Time
		due = &d
	}

	var deleted *time.Time
	if dbCard.DeletedAt.Valid {
		d := dbCard.DeletedAt.Time
		deleted = &d
	}

	var startDate *time.Time
	if dbCard.StartDate.Valid {
		s := dbCard.StartDate.Time
		startDate = &s
	}

	var completionDate *time.Time
	if dbCard.CompletionDate.Valid {
		c := dbCard.CompletionDate.Time
		completionDate = &c
	}

	var lastActivityAt *time.Time
	if dbCard.LastActivityAt.Valid {
		l := dbCard.LastActivityAt.Time
		lastActivityAt = &l
	}

	pos := 0.0
	if dbCard.Position.Big != nil {
		f, _ := dbCard.Position.Big.Float64()
		pos = f
	}

	var estimatedHours *float64
	if dbCard.EstimatedHours.Big != nil {
		f, _ := dbCard.EstimatedHours.Big.Float64()
		estimatedHours = &f
	}

	var actualHours *float64
	if dbCard.ActualHours.Big != nil {
		f, _ := dbCard.ActualHours.Big.Float64()
		actualHours = &f
	}

	return Card{
		ID:             dbCard.ID,
		ListID:         dbCard.ListID,
		Name:           dbCard.Name,
		Alias:          dbCard.Alias.String,
		Description:    desc,
		Position:       pos,
		DueDate:        due,
		Priority:       CardPriority(dbCard.Priority),
		Labels:         labels,
		IsArchived:     dbCard.IsArchived,
		CreatedBy:      dbCard.CreatedBy.Ptr(),
		CreatedAt:      dbCard.CreatedAt,
		UpdatedAt:      dbCard.UpdatedAt,
		DeletedAt:      deleted,
		AssignedTo:     dbCard.AssignedTo.Ptr(),
		Attachments:    attachments,
		EstimatedHours: estimatedHours,
		ActualHours:    actualHours,
		StartDate:      startDate,
		CompletionDate: completionDate,
		Tags:           dbCard.Tags,
		Checklist:      checklist,
		LastActivityAt: lastActivityAt,
		UpdatedBy:      dbCard.UpdatedBy.Ptr(),
	}
}
