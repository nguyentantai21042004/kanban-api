package repository

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/role"
	pag "gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type GetOneOptions struct {
	Filter role.Filter
}

type GetOptions struct {
	Filter   role.Filter
	PagQuery pag.PaginateQuery
}

type ListOptions struct {
	Filter role.Filter
}
