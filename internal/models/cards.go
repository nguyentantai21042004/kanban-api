package models

import (
	"encoding/json"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
)

type Card struct {
	ID          string       `json:"id"`
	ListID      string       `json:"list_id"`
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Position    int          `json:"position"`
	DueDate     *time.Time   `json:"due_date,omitempty"`
	Priority    CardPriority `json:"priority"`
	Labels      []string     `json:"labels,omitempty"`
	IsArchived  bool         `json:"is_archived"`
	CreatedBy   *string      `json:"created_by,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   *time.Time   `json:"deleted_at,omitempty"`
}

type CardPriority string

const (
	CardPriorityLow    CardPriority = "low"
	CardPriorityMedium CardPriority = "medium"
	CardPriorityHigh   CardPriority = "high"
)

func NewCard(dbCard dbmodels.Card) Card {
	labels := []string{}
	if dbCard.Labels.Valid {
		_ = json.Unmarshal(dbCard.Labels.JSON, &labels)
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
	pos := 0
	if dbCard.Position.Big != nil {
		f, _ := dbCard.Position.Big.Float64()
		pos = int(f)
	}
	return Card{
		ID:          dbCard.ID,
		ListID:      dbCard.ListID,
		Title:       dbCard.Title,
		Description: desc,
		Position:    pos,
		DueDate:     due,
		Priority:    CardPriority(dbCard.Priority),
		Labels:      labels,
		IsArchived:  dbCard.IsArchived,
		CreatedBy:   dbCard.CreatedBy.Ptr(),
		CreatedAt:   dbCard.CreatedAt,
		UpdatedAt:   dbCard.UpdatedAt,
		DeletedAt:   deleted,
	}
}
