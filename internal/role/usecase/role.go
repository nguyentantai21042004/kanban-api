package usecase

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/role"
)

func (uc *usecase) Detail(ctx context.Context, sc models.Scope, ID string) (role.DetailOutput, error) {
	roleModel, err := uc.repo.Detail(ctx, sc, ID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.role.usecase.Detail.repo.Detail: %v", err)
		return role.DetailOutput{}, err
	}

	return role.DetailOutput{Role: roleModel}, nil
}

func (uc *usecase) GetOne(ctx context.Context, sc models.Scope, ip role.GetOneInput) (role.GetOneOutput, error) {
	roleModel, err := uc.repo.GetOne(ctx, sc, role.GetOneOptions{Filter: ip.Filter})
	if err != nil {
		uc.l.Errorf(ctx, "internal.role.usecase.GetOne.repo.GetOne: %v", err)
		return role.GetOneOutput{}, err
	}

	return role.GetOneOutput{Role: roleModel}, nil
}

func (uc *usecase) Get(ctx context.Context, sc models.Scope, ip role.GetInput) (role.GetOutput, error) {
	roles, paginator, err := uc.repo.Get(ctx, sc, role.GetOptions{
		Filter:   ip.Filter,
		PagQuery: ip.PagQuery,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.role.usecase.Get.repo.Get: %v", err)
		return role.GetOutput{}, err
	}

	return role.GetOutput{
		Roles:     roles,
		Paginator: paginator,
	}, nil
}

func (uc *usecase) List(ctx context.Context, sc models.Scope, ip role.ListInput) (role.ListOutput, error) {
	roles, err := uc.repo.List(ctx, sc, role.ListOptions{Filter: ip.Filter})
	if err != nil {
		uc.l.Errorf(ctx, "internal.role.usecase.List.repo.List: %v", err)
		return role.ListOutput{}, err
	}

	return role.ListOutput{Roles: roles}, nil
}
