package models

import (
	"time"
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
