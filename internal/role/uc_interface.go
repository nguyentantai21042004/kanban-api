package role

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

//go:generate mockery --name UseCase
type UseCase interface {
	Detail(ctx context.Context, sc models.Scope, ID string) (DetailOutput, error)
	GetOne(ctx context.Context, sc models.Scope, ip GetOneInput) (GetOneOutput, error)
	Get(ctx context.Context, sc models.Scope, ip GetInput) (GetOutput, error)
	List(ctx context.Context, sc models.Scope, ip ListInput) (ListOutput, error)
}
