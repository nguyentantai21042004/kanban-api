package models

import (
	"time"

	"github.com/aarondl/sqlboiler/v4/types"
)

type List struct {
	ID         string        `json:"id"`
	BoardID    string        `json:"board_id"`
	Title      string        `json:"title"`
	Position   types.Decimal `json:"position"`
	IsArchived bool          `json:"is_archived"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	DeletedAt  *time.Time    `json:"deleted_at,omitempty"`
}
