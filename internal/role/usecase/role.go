package usecase

import (
	"context"

	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/internal/role"
	"github.com/nguyentantai21042004/kanban-api/internal/role/repository"
)

func (uc *usecase) Detail(ctx context.Context, sc models.Scope, ID string) (models.Role, error) {
	r, err := uc.repo.Detail(ctx, sc, ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.role.usecase.Detail.repo.Detail: %v", err)
			return models.Role{}, role.ErrRoleNotFound
		}
		uc.l.Errorf(ctx, "internal.role.usecase.Detail.repo.Detail: %v", err)
		return models.Role{}, err
	}

	return r, nil
}

func (uc *usecase) List(ctx context.Context, sc models.Scope, ip role.ListInput) ([]models.Role, error) {
	rls, err := uc.repo.List(ctx, sc, repository.ListOptions{Filter: ip.Filter})
	if err != nil {
		uc.l.Errorf(ctx, "internal.role.usecase.List.repo.List: %v", err)
		return []models.Role{}, err
	}

	return rls, nil
}
