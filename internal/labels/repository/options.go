package repository

import (
	"github.com/nguyentantai21042004/kanban-api/internal/labels"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/pkg/paginator"
)

type GetOptions struct {
	Filter   labels.Filter
	PagQuery paginator.PaginateQuery
}

type CreateOptions struct {
	BoardID string
	Name    string
	Color   string
}

type UpdateOptions struct {
	ID       string
	Name     string
	Color    string
	OldModel models.Label
}
