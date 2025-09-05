package usecase

import (
	"context"

	"github.com/nguyentantai21042004/kanban-api/internal/admin"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/internal/role"
)

func (uc implUsecase) Roles(ctx context.Context, sc models.Scope) ([]admin.RoleItem, error) {
	roles, err := uc.roleUC.List(ctx, sc, role.ListInput{})
	if err != nil {
		uc.l.Errorf(ctx, "internal.admin.usecase.Roles.roleUC.List: %v", err)
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
