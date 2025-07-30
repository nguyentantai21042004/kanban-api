package models

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
)

type Board struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Alias       string     `json:"alias"`
	Description *string    `json:"description,omitempty"`
	CreatedBy   *string    `json:"created_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func NewBoard(dbBoard dbmodels.Board) Board {
	board := Board{
		ID:          dbBoard.ID,
		Name:        dbBoard.Name,
		Alias:       dbBoard.Alias.String,
		Description: dbBoard.Description.Ptr(),
		CreatedBy:   dbBoard.CreatedBy.Ptr(),
		CreatedAt:   dbBoard.CreatedAt,
		UpdatedAt:   dbBoard.UpdatedAt,
		DeletedAt:   dbBoard.DeletedAt.Ptr(),
	}

	return board
}
