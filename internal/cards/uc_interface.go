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

	// New methods for enhanced functionality
	Assign(ctx context.Context, sc models.Scope, ip AssignInput) (DetailOutput, error)
	Unassign(ctx context.Context, sc models.Scope, ip UnassignInput) (DetailOutput, error)
	AddAttachment(ctx context.Context, sc models.Scope, ip AddAttachmentInput) (DetailOutput, error)
	RemoveAttachment(ctx context.Context, sc models.Scope, ip RemoveAttachmentInput) (DetailOutput, error)
	UpdateTimeTracking(ctx context.Context, sc models.Scope, ip UpdateTimeTrackingInput) (DetailOutput, error)
	UpdateChecklist(ctx context.Context, sc models.Scope, ip UpdateChecklistInput) (DetailOutput, error)
	AddTag(ctx context.Context, sc models.Scope, ip AddTagInput) (DetailOutput, error)
	RemoveTag(ctx context.Context, sc models.Scope, ip RemoveTagInput) (DetailOutput, error)
	SetStartDate(ctx context.Context, sc models.Scope, ip SetStartDateInput) (DetailOutput, error)
	SetCompletionDate(ctx context.Context, sc models.Scope, ip SetCompletionDateInput) (DetailOutput, error)
}
