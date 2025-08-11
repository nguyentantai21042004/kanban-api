package models

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
)

type Label struct {
	ID        string     `json:"id"`
	BoardID   string     `json:"board_id"`
	Name      string     `json:"name"`
	Color     string     `json:"color"`
	CreatedBy *string    `json:"created_by,omitempty"`
	UpdatedBy *string    `json:"updated_by,omitempty"`
	DeletedBy *string    `json:"deleted_by,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func NewLabel(dbLabel dbmodels.Label) Label {
	return Label{
		ID:        dbLabel.ID,
		BoardID:   dbLabel.BoardID,
		Name:      dbLabel.Name,
		Color:     dbLabel.Color,
		CreatedBy: dbLabel.CreatedBy.Ptr(),
		UpdatedBy: dbLabel.UpdatedBy.Ptr(),
		DeletedBy: dbLabel.DeletedBy.Ptr(),
		CreatedAt: dbLabel.CreatedAt,
		UpdatedAt: dbLabel.UpdatedAt,
		DeletedAt: dbLabel.DeletedAt.Ptr(),
	}
}
