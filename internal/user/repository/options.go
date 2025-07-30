package repository

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
)

type CreateOptions struct {
	User models.User
}

type UpdateOptions struct {
	User models.User
}

type GetOneOptions struct {
	Username string
}

type ListOptions struct {
	Filter user.Filter
}
