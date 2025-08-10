package models

import (
	"encoding/json"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
)

type CardActivity struct {
	ID         string                 `json:"id"`
	CardID     string                 `json:"card_id"`
	ActionType CardActionType         `json:"action_type"`
	OldData    map[string]interface{} `json:"old_data,omitempty"`
	NewData    map[string]interface{} `json:"new_data,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
	DeletedAt  *time.Time             `json:"deleted_at,omitempty"`
}

type CardActionType string

const (
	CardActionTypeCreated   CardActionType = "created"
	CardActionTypeUpdated   CardActionType = "updated"
	CardActionTypeDeleted   CardActionType = "deleted"
	CardActionTypeMoved     CardActionType = "moved"
	CardActionTypeAssigned  CardActionType = "assigned"
	CardActionTypeCommented CardActionType = "commented"
)

func NewCardActivity(dbCardActivity dbmodels.CardActivity) CardActivity {
	var oldData map[string]interface{}
	if dbCardActivity.OldData.Valid {
		_ = json.Unmarshal(dbCardActivity.OldData.JSON, &oldData)
	}

	var newData map[string]interface{}
	if dbCardActivity.NewData.Valid {
		_ = json.Unmarshal(dbCardActivity.NewData.JSON, &newData)
	}

	return CardActivity{
		ID:         dbCardActivity.ID,
		CardID:     dbCardActivity.CardID,
		ActionType: CardActionType(dbCardActivity.ActionType),
		OldData:    oldData,
		NewData:    newData,
		CreatedAt:  dbCardActivity.CreatedAt,
		UpdatedAt:  dbCardActivity.UpdatedAt,
		DeletedAt:  dbCardActivity.DeletedAt.Ptr(),
	}
}
