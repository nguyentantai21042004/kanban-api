package admin

import (
	"context"

	"github.com/nguyentantai21042004/kanban-api/internal/models"
)

//go:generate mockery --name UseCase
type UseCase interface {
	Dashboard(ctx context.Context, sc models.Scope, ip DashboardInput) (DashboardOutput, error)
	Users(ctx context.Context, sc models.Scope, ip UsersInput) (UsersOutput, error)
	CreateUser(ctx context.Context, sc models.Scope, ip CreateUserInput) (UserItem, error)
	UpdateUser(ctx context.Context, sc models.Scope, id string, ip UpdateUserInput) (UserItem, error)
	Roles(ctx context.Context, sc models.Scope) ([]RoleItem, error)
	Health(ctx context.Context, sc models.Scope) (HealthOutput, error)
}
