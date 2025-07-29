package models

import (
	"encoding/json"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
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

func NewCardActivity(dbAct dbmodels.CardActivity) CardActivity {
	var oldData map[string]any
	if dbAct.OldData.Valid && len(dbAct.OldData.JSON) > 0 {
		_ = json.Unmarshal(dbAct.OldData.JSON, &oldData)
	}
	var newData map[string]any
	if dbAct.NewData.Valid && len(dbAct.NewData.JSON) > 0 {
		_ = json.Unmarshal(dbAct.NewData.JSON, &newData)
	}
	var deleted *time.Time
	if dbAct.DeletedAt.Valid {
		d := dbAct.DeletedAt.Time
		deleted = &d
	}
	return CardActivity{
		ID:         dbAct.ID,
		CardID:     dbAct.CardID,
		ActionType: CardActionType(dbAct.ActionType),
		OldData:    oldData,
		NewData:    newData,
		CreatedAt:  dbAct.CreatedAt,
		UpdatedAt:  dbAct.UpdatedAt,
		DeletedAt:  deleted,
	}
}
