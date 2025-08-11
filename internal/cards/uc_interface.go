package cards

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

//go:generate mockery --name UseCase
type UseCase interface {
	CoreUseCase
	EnhancedUseCase
}

type CoreUseCase interface {
	Detail(ctx context.Context, sc models.Scope, ID string) (DetailOutput, error)
	Get(ctx context.Context, sc models.Scope, ip GetInput) (GetOutput, error)
	Move(ctx context.Context, sc models.Scope, ip MoveInput) error
	Create(ctx context.Context, sc models.Scope, ip CreateInput) (DetailOutput, error)
	Update(ctx context.Context, sc models.Scope, ip UpdateInput) (DetailOutput, error)
	Delete(ctx context.Context, sc models.Scope, ids []string) error
}

type EnhancedUseCase interface {
	GetActivities(ctx context.Context, sc models.Scope, ip GetActivitiesInput) (GetActivitiesOutput, error)
	Assign(ctx context.Context, sc models.Scope, ip AssignInput) error
	Unassign(ctx context.Context, sc models.Scope, ip UnassignInput) error
	AddAttachment(ctx context.Context, sc models.Scope, ip AddAttachmentInput) error
	RemoveAttachment(ctx context.Context, sc models.Scope, ip RemoveAttachmentInput) error
	UpdateTimeTracking(ctx context.Context, sc models.Scope, ip UpdateTimeTrackingInput) error
	UpdateChecklist(ctx context.Context, sc models.Scope, ip UpdateChecklistInput) error
	AddTag(ctx context.Context, sc models.Scope, ip AddTagInput) error
	RemoveTag(ctx context.Context, sc models.Scope, ip RemoveTagInput) error
	SetStartDate(ctx context.Context, sc models.Scope, ip SetStartDateInput) error
	SetCompletionDate(ctx context.Context, sc models.Scope, ip SetCompletionDateInput) error
}
