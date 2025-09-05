package models

import (
	"time"

	"github.com/nguyentantai21042004/kanban-api/internal/dbmodels"
)

type Role struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	Alias       string     `json:"alias"`
	Description string     `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func NewRole(r dbmodels.Role) Role {
	return Role{
		ID:          r.ID,
		Name:        r.Name,
		Code:        r.Code,
		Alias:       r.Alias,
		Description: r.Description.String,
		CreatedAt:   r.CreatedAt.Time,
		UpdatedAt:   r.UpdatedAt.Time,
		DeletedAt:   r.DeletedAt.Ptr(),
	}
}
