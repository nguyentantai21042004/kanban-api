package upload

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

//go:generate mockery --name UseCase
type UseCase interface {
	Create(ctx context.Context, sc models.Scope, ip CreateInput) (UploadOutput, error)
	Detail(ctx context.Context, sc models.Scope, ID string) (UploadOutput, error)
	Get(ctx context.Context, sc models.Scope, ip GetInput) (GetOutput, error)
}
