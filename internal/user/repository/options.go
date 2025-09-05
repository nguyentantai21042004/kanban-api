package repository

import (
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/internal/user"
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
