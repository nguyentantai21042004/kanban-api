package usecase

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/boards"
	"gitlab.com/tantai-kanban/kanban-api/internal/boards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

func (uc implUsecase) Get(ctx context.Context, sc models.Scope, ip boards.GetInput) (boards.GetOutput, error) {
	u, err := uc.userUC.DetailMe(ctx, sc)
	if err != nil {
		uc.l.Errorf(ctx, "internal.boards.usecase.Get.userUC.Detail: %v", err)
		return boards.GetOutput{}, err
	}

	rl, err := uc.roleUC.Detail(ctx, sc, u.User.RoleID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.boards.usecase.Get.roleUC.Detail: %v", err)
		return boards.GetOutput{}, err
	}

	// Only admin can see all boards
	if rl.Code != models.ADMIN_ROLE {
		ip.Filter.CreatedBy = u.User.ID
	}

	b, p, err := uc.repo.Get(ctx, sc, repository.GetOptions{
		Filter: boards.Filter{
			IDs:       ip.Filter.IDs,
			Keyword:   util.BuildAlias(ip.Filter.Keyword),
			CreatedBy: ip.Filter.CreatedBy,
		},
		PagQuery: ip.PagQuery,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.boards.usecase.Get.repo.Get: %v", err)
		return boards.GetOutput{}, err
	}

	uIDs := make([]string, len(b))
	for i, b := range b {
		uIDs[i] = *b.CreatedBy
	}
	uIDs = util.RemoveDuplicates(uIDs)
	us, err := uc.userUC.List(ctx, sc, user.ListInput{
		Filter: user.Filter{
			IDs: uIDs,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.boards.usecase.Get.userUC.List: %v", err)
		return boards.GetOutput{}, err
	}

	return boards.GetOutput{
		Boards:     b,
		Users:      us,
		Pagination: p,
	}, nil
}

func (uc implUsecase) Create(ctx context.Context, sc models.Scope, ip boards.CreateInput) (boards.DetailOutput, error) {
	b, err := uc.repo.Create(ctx, sc, repository.CreateOptions{
		Name:        ip.Name,
		Alias:       util.BuildAlias(ip.Name),
		Description: ip.Description,
	})

	if err != nil {
		uc.l.Errorf(ctx, "internal.boards.usecase.Create.repo.Create: %v", err)
		return boards.DetailOutput{}, err
	}

	u, err := uc.userUC.Detail(ctx, sc, *b.CreatedBy)
	if err != nil {
		if err == user.ErrUserNotFound {
			uc.l.Warnf(ctx, "internal.boards.usecase.Create.userUC.DetailMe: %v", err)
			return boards.DetailOutput{}, err
		}
		uc.l.Errorf(ctx, "internal.boards.usecase.Create.userUC.Detail: %v", err)
		return boards.DetailOutput{}, err
	}

	return boards.DetailOutput{
		Board: b,
		User:  u.User,
	}, nil
}

func (uc implUsecase) Update(ctx context.Context, sc models.Scope, ip boards.UpdateInput) (boards.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.boards.usecase.Update.repo.Detail.NotFound: %v", err)
			return boards.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.boards.usecase.Update.repo.Detail: %v", err)
		return boards.DetailOutput{}, err
	}

	b, err := uc.repo.Update(ctx, sc, repository.UpdateOptions{
		ID:          ip.ID,
		Name:        ip.Name,
		Description: ip.Description,
		OldModel:    oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.boards.usecase.Update.repo.Update: %v", err)
		return boards.DetailOutput{}, err
	}

	return boards.DetailOutput{
		Board: b,
	}, nil
}

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, ID string) (boards.DetailOutput, error) {
	b, err := uc.repo.Detail(ctx, sc, ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.boards.usecase.Detail.repo.Detail.NotFound: %v", err)
			return boards.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.boards.usecase.Detail.repo.Detail: %v", err)
		return boards.DetailOutput{}, err
	}
	return boards.DetailOutput{
		Board: b,
	}, nil
}

func (uc implUsecase) Delete(ctx context.Context, sc models.Scope, ids []string) error {
	if len(ids) == 0 {
		uc.l.Warnf(ctx, "internal.boards.usecase.Delete.ids.Empty")
		return boards.ErrFieldRequired
	}

	err := uc.repo.Delete(ctx, sc, ids)
	if err != nil {
		uc.l.Errorf(ctx, "internal.boards.usecase.Delete.repo.Delete: %v", err)
		return err
	}
	return nil
}
