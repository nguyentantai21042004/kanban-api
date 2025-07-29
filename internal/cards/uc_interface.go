package cards

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

//go:generate mockery --name UseCase
type UseCase interface {
	Get(ctx context.Context, sc models.Scope, ip GetInput) (GetOutput, error)
	Create(ctx context.Context, sc models.Scope, ip CreateInput) (DetailOutput, error)
	Update(ctx context.Context, sc models.Scope, ip UpdateInput) (DetailOutput, error)
	Move(ctx context.Context, sc models.Scope, ip MoveInput) (DetailOutput, error)
	Detail(ctx context.Context, sc models.Scope, ID string) (DetailOutput, error)
	Delete(ctx context.Context, sc models.Scope, ids []string) error
	GetActivities(ctx context.Context, sc models.Scope, ip GetActivitiesInput) (GetActivitiesOutput, error)
}
