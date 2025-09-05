package repository

import (
	"github.com/nguyentantai21042004/kanban-api/internal/boards"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/pkg/paginator"
)

type ListOptions struct {
	Filter boards.Filter
}

type GetOptions struct {
	Filter   boards.Filter
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
	OldModel    models.Board
}
