package lists

import (
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/pkg/paginator"
)

type Filter struct {
	IDs       []string
	BoardID   string
	Keyword   string
	CreatedBy string
}

type GetInput struct {
	Filter   Filter
	PagQuery paginator.PaginateQuery
}

type CreateInput struct {
	BoardID string
	Name    string
}

type UpdateInput struct {
	ID   string
	Name string
}

type MoveInput struct {
	ID       string
	BoardID  string
	AfterID  string
	BeforeID string
}

type GetOutput struct {
	Lists      []models.List
	Pagination paginator.Paginator
}

type DetailOutput struct {
	Board models.Board
	List  models.List
}
