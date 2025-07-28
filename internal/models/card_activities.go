package models

import (
	"time"
)

type CardActivity struct {
	ID         string         `json:"id"`
	CardID     string         `json:"card_id"`
	ActionType CardActionType `json:"action_type"`
	OldData    map[string]any `json:"old_data,omitempty"`
	NewData    map[string]any `json:"new_data,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  *time.Time     `json:"deleted_at,omitempty"`
}

type CardActionType string

const (
	CardActionTypeCreated   CardActionType = "created"
	CardActionTypeMoved     CardActionType = "moved"
	CardActionTypeUpdated   CardActionType = "updated"
	CardActionTypeCommented CardActionType = "commented"
)
