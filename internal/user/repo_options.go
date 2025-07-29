package user

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

type CreateOptions struct {
	User models.User
}

type UpdateOptions struct {
	User models.User
}
