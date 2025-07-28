package usecase

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/labels"
	"gitlab.com/tantai-kanban/kanban-api/internal/labels/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

func (uc implUsecase) Get(ctx context.Context, sc models.Scope, ip labels.GetInput) (labels.GetOutput, error) {
	u, p, err := uc.repo.Get(ctx, sc, repository.GetOptions{
		Filter:   ip.Filter,
		PagQuery: ip.PagQuery,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.labels.usecase.Get.repo.Get: %v", err)
		return labels.GetOutput{}, err
	}

	return labels.GetOutput{
		Labels:     u,
		Pagination: p,
	}, nil
}

func (uc implUsecase) Create(ctx context.Context, sc models.Scope, ip labels.CreateInput) (labels.DetailOutput, error) {
	b, err := uc.repo.Create(ctx, sc, repository.CreateOptions{
		BoardID: ip.BoardID,
		Name:    ip.Name,
		Color:   ip.Color,
	})

	if err != nil {
		uc.l.Errorf(ctx, "internal.labels.usecase.Create.repo.Create: %v", err)
		return labels.DetailOutput{}, err
	}

	return labels.DetailOutput{
		Label: b,
	}, nil
}

func (uc implUsecase) Update(ctx context.Context, sc models.Scope, ip labels.UpdateInput) (labels.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.labels.usecase.Update.repo.Detail.NotFound: %v", err)
			return labels.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.labels.usecase.Update.repo.Detail: %v", err)
		return labels.DetailOutput{}, err
	}

	b, err := uc.repo.Update(ctx, sc, repository.UpdateOptions{
		ID:       ip.ID,
		Name:     ip.Name,
		Color:    ip.Color,
		OldModel: oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.labels.usecase.Update.repo.Update: %v", err)
		return labels.DetailOutput{}, err
	}

	return labels.DetailOutput{
		Label: b,
	}, nil
}

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, ID string) (labels.DetailOutput, error) {
	b, err := uc.repo.Detail(ctx, sc, ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.labels.usecase.Detail.repo.Detail.NotFound: %v", err)
			return labels.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.labels.usecase.Detail.repo.Detail: %v", err)
		return labels.DetailOutput{}, err
	}
	return labels.DetailOutput{
		Label: b,
	}, nil
}

func (uc implUsecase) Delete(ctx context.Context, sc models.Scope, ids []string) error {
	if len(ids) == 0 {
		uc.l.Warnf(ctx, "internal.labels.usecase.Delete.ids.Empty")
		return labels.ErrFieldRequired
	}

	err := uc.repo.Delete(ctx, sc, ids)
	if err != nil {
		uc.l.Errorf(ctx, "internal.labels.usecase.Delete.repo.Delete: %v", err)
		return err
	}
	return nil
}
