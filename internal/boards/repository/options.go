package repository

import "gitlab.com/tantai-kanban/kanban-api/pkg/paginator"

type Filter struct {
	IDs     []string
	Keyword string
}

type GetOptions struct {
	Filter   Filter
	PagQuery paginator.PaginateQuery
}

type CreateOptions struct {
	Name        string
	Description string
	Alias       string
}

type UpdateOptions struct {
	ID          string
	Name        string
	Description string
	Alias       string
}
