package models

import (
	"time"

	"github.com/aarondl/sqlboiler/v4/types"
	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
)

type List struct {
	ID         string        `json:"id"`
	BoardID    string        `json:"board_id"`
	Name       string        `json:"name"`
	Position   types.Decimal `json:"position"`
	IsArchived bool          `json:"is_archived"`
	CreatedBy  *string       `json:"created_by,omitempty"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	DeletedAt  *time.Time    `json:"deleted_at,omitempty"`
}

func NewList(dbList dbmodels.List) List {
	return List{
		ID:         dbList.ID,
		BoardID:    dbList.BoardID,
		Name:       dbList.Name,
		Position:   dbList.Position,
		IsArchived: dbList.IsArchived,
		CreatedBy:  dbList.CreatedBy.Ptr(),
		CreatedAt:  dbList.CreatedAt,
		UpdatedAt:  dbList.UpdatedAt,
		DeletedAt:  dbList.DeletedAt.Ptr(),
	}
}
