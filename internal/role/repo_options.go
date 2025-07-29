package role

import pag "gitlab.com/tantai-kanban/kanban-api/pkg/paginator"

type GetOneOptions struct {
	Filter Filter
}

type GetOptions struct {
	Filter   Filter
	PagQuery pag.PaginateQuery
}

type ListOptions struct {
	Filter Filter
}
