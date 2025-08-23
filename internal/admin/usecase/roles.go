package usecase

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/admin"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/role"
)

func (uc implUsecase) Roles(ctx context.Context, sc models.Scope) ([]admin.RoleItem, error) {
	roles, err := uc.roleUC.List(ctx, sc, role.ListInput{})
	if err != nil {
		return nil, err
	}
	
	items := make([]admin.RoleItem, len(roles))
	for i, r := range roles {
		items[i] = admin.RoleItem{
			ID:    r.ID,
			Name:  r.Name,
			Alias: r.Alias,
		}
	}
	
	return items, nil
}