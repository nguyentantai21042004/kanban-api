package role

import (
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	pag "github.com/nguyentantai21042004/kanban-api/pkg/paginator"
)

type Filter struct {
	IDs      []string `json:"ids"`
	Code     string   `json:"code"`
	IsActive bool     `json:"is_active"`
}

type GetOneInput struct {
	Filter Filter
}

type GetInput struct {
	Filter   Filter
	PagQuery pag.PaginateQuery
}

type ListInput struct {
	Filter Filter
}

type DetailOutput struct {
	Role models.Role
}

type GetOneOutput struct {
	Role models.Role
}

type GetOutput struct {
	Roles     []models.Role
	Paginator pag.Paginator
}

type ListOutput struct {
	Roles []models.Role
}
