package usecase

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/lists"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

func (uc implUsecase) Get(ctx context.Context, sc models.Scope, ip lists.GetInput) (lists.GetOutput, error) {
	u, p, err := uc.repo.Get(ctx, sc, repository.GetOptions{
		Filter:   ip.Filter,
		PagQuery: ip.PagQuery,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Get.repo.Get: %v", err)
		return lists.GetOutput{}, err
	}

	return lists.GetOutput{
		Lists:      u,
		Pagination: p,
	}, nil
}

func (uc implUsecase) Create(ctx context.Context, sc models.Scope, ip lists.CreateInput) (lists.DetailOutput, error) {
	b, err := uc.repo.Create(ctx, sc, repository.CreateOptions{
		BoardID:  ip.BoardID,
		Title:    ip.Title,
		Position: ip.Position,
	})

	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Create.repo.Create: %v", err)
		return lists.DetailOutput{}, err
	}

	return lists.DetailOutput{
		List: b,
	}, nil
}

func (uc implUsecase) Update(ctx context.Context, sc models.Scope, ip lists.UpdateInput) (lists.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.lists.usecase.Update.repo.Detail.NotFound: %v", err)
			return lists.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.lists.usecase.Update.repo.Detail: %v", err)
		return lists.DetailOutput{}, err
	}

	b, err := uc.repo.Update(ctx, sc, repository.UpdateOptions{
		ID:       ip.ID,
		Title:    ip.Title,
		Position: ip.Position,
		OldModel: oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Update.repo.Update: %v", err)
		return lists.DetailOutput{}, err
	}

	return lists.DetailOutput{
		List: b,
	}, nil
}

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, ID string) (lists.DetailOutput, error) {
	b, err := uc.repo.Detail(ctx, sc, ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.lists.usecase.Detail.repo.Detail.NotFound: %v", err)
			return lists.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.lists.usecase.Detail.repo.Detail: %v", err)
		return lists.DetailOutput{}, err
	}
	return lists.DetailOutput{
		List: b,
	}, nil
}

func (uc implUsecase) Delete(ctx context.Context, sc models.Scope, ids []string) error {
	if len(ids) == 0 {
		uc.l.Warnf(ctx, "internal.lists.usecase.Delete.ids.Empty")
		return lists.ErrFieldRequired
	}

	err := uc.repo.Delete(ctx, sc, ids)
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Delete.repo.Delete: %v", err)
		return err
	}
	return nil
}
