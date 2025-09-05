package http

import (
	"github.com/nguyentantai21042004/kanban-api/internal/models"
)

type roleItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Alias       string `json:"alias"`
	Description string `json:"description"`
}

func (h handler) newItem(o models.Role) roleItem {
	return roleItem{
		ID:          o.ID,
		Name:        o.Name,
		Code:        o.Code,
		Alias:       o.Alias,
		Description: o.Description,
	}
}
