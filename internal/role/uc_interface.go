package role

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

//go:generate mockery --name UseCase
type UseCase interface {
	Detail(ctx context.Context, sc models.Scope, ID string) (models.Role, error)
	List(ctx context.Context, sc models.Scope, ip ListInput) ([]models.Role, error)
}
