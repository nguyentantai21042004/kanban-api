package comments

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type Filter struct {
	IDs      []string
	Keyword  string
	CardID   string
	UserID   string
	ParentID string
}

type GetInput struct {
	Filter   Filter
	PagQuery paginator.PaginateQuery
}

type CreateInput struct {
	CardID   string
	Content  string
	ParentID *string
}

type UpdateInput struct {
	ID      string
	Content string
}

type GetOutput struct {
	Comments   []models.Comment
	Users      []models.User
	Pagination paginator.Paginator
}

type DetailOutput struct {
	Comment models.Comment
	User    models.User
}

type CommentWithDetailsOutput struct {
	Comment models.Comment
	User    models.User
	Card    models.Card
	Parent  *models.Comment
	Replies []models.Comment
}
