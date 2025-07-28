package models

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
)

type Label struct {
	ID      string `json:"id"`
	BoardID string `json:"board_id"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

func NewLabel(dbLabel dbmodels.Label) Label {
	return Label{
		ID:      dbLabel.ID,
		BoardID: dbLabel.BoardID,
		Name:    dbLabel.Name,
		Color:   dbLabel.Color,
	}
}
