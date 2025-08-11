package repository

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

//go:generate mockery --name Repository
type Repository interface {
	CoreRepository
	EnhancedRepository
}

type CoreRepository interface {
	Detail(ctx context.Context, sc models.Scope, id string) (models.Card, error)
	List(ctx context.Context, sc models.Scope, opts ListOptions) ([]models.Card, error)
	Get(ctx context.Context, sc models.Scope, opts GetOptions) ([]models.Card, paginator.Paginator, error)
	Move(ctx context.Context, sc models.Scope, opts MoveOptions) (models.Card, error)
	Create(ctx context.Context, sc models.Scope, opts CreateOptions) (models.Card, error)
	Update(ctx context.Context, sc models.Scope, opts UpdateOptions) (models.Card, error)
	Delete(ctx context.Context, sc models.Scope, ids []string) error
}

type EnhancedRepository interface {
	GetPosition(ctx context.Context, sc models.Scope, opts GetPositionOptions) (string, error)
	GetActivities(ctx context.Context, sc models.Scope, opts GetActivitiesOptions) ([]models.CardActivity, paginator.Paginator, error)
	Assign(ctx context.Context, sc models.Scope, opts AssignOptions) (models.Card, error)
	Unassign(ctx context.Context, sc models.Scope, opts UnassignOptions) (models.Card, error)
	AddAttachment(ctx context.Context, sc models.Scope, opts AddAttachmentOptions) (models.Card, error)
	RemoveAttachment(ctx context.Context, sc models.Scope, opts RemoveAttachmentOptions) (models.Card, error)
	UpdateTimeTracking(ctx context.Context, sc models.Scope, opts UpdateTimeTrackingOptions) (models.Card, error)
	UpdateChecklist(ctx context.Context, sc models.Scope, opts UpdateChecklistOptions) (models.Card, error)
	AddTag(ctx context.Context, sc models.Scope, opts AddTagOptions) (models.Card, error)
	RemoveTag(ctx context.Context, sc models.Scope, opts RemoveTagOptions) (models.Card, error)
	SetStartDate(ctx context.Context, sc models.Scope, opts SetStartDateOptions) (models.Card, error)
	SetCompletionDate(ctx context.Context, sc models.Scope, opts SetCompletionDateOptions) (models.Card, error)
}
