package repository

import (
	"github.com/nguyentantai21042004/kanban-api/internal/role"
	pag "github.com/nguyentantai21042004/kanban-api/pkg/paginator"
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
