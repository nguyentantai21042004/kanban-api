package usecase

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/cards"
	"gitlab.com/tantai-kanban/kanban-api/internal/comments"
	"gitlab.com/tantai-kanban/kanban-api/internal/comments/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

// broadcastCommentEvent broadcasts comment events to WebSocket clients
func (uc implUsecase) broadcastCommentEvent(ctx context.Context, cardID, eventType string, data interface{}, userID string) {
	if uc.wsHub == nil {
		return
	}

	uc.wsHub.BroadcastToBoard(cardID, eventType, data, userID)
}

func (uc implUsecase) Get(ctx context.Context, sc models.Scope, ip comments.GetInput) (comments.GetOutput, error) {
	_, err := uc.userUC.DetailMe(ctx, sc)
	if err != nil {
		uc.l.Errorf(ctx, "internal.comments.usecase.Get.userUC.Detail: %v", err)
		return comments.GetOutput{}, err
	}

	c, p, err := uc.repo.Get(ctx, sc, repository.GetOptions{
		Filter: comments.Filter{
			IDs:      ip.Filter.IDs,
			Keyword:  ip.Filter.Keyword,
			CardID:   ip.Filter.CardID,
			UserID:   ip.Filter.UserID,
			ParentID: ip.Filter.ParentID,
		},
		PagQuery: ip.PagQuery,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.comments.usecase.Get.repo.Get: %v", err)
		return comments.GetOutput{}, err
	}

	uIDs := make([]string, len(c))
	for i, comment := range c {
		uIDs[i] = comment.UserID
	}
	uIDs = util.RemoveDuplicates(uIDs)
	us, err := uc.userUC.List(ctx, sc, user.ListInput{
		Filter: user.Filter{
			IDs: uIDs,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.comments.usecase.Get.userUC.List: %v", err)
		return comments.GetOutput{}, err
	}

	return comments.GetOutput{
		Comments:   c,
		Users:      us,
		Pagination: p,
	}, nil
}

func (uc implUsecase) Create(ctx context.Context, sc models.Scope, ip comments.CreateInput) (comments.DetailOutput, error) {
	// Verify card exists
	_, err := uc.cardsUC.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == cards.ErrCardNotFound {
			uc.l.Warnf(ctx, "internal.comments.usecase.Create.cardsUC.Detail.CardNotFound: %v", err)
			return comments.DetailOutput{}, comments.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.comments.usecase.Create.cardsUC.Detail: %v", err)
		return comments.DetailOutput{}, err
	}

	// Verify parent comment exists if provided
	if ip.ParentID != nil {
		_, err := uc.repo.Detail(ctx, sc, *ip.ParentID)
		if err != nil {
			if err == repository.ErrNotFound {
				uc.l.Warnf(ctx, "internal.comments.usecase.Create.repo.Detail.ParentNotFound: %v", err)
				return comments.DetailOutput{}, comments.ErrParentCommentNotFound
			}
			uc.l.Errorf(ctx, "internal.comments.usecase.Create.repo.Detail: %v", err)
			return comments.DetailOutput{}, err
		}
	}

	c, err := uc.repo.Create(ctx, sc, repository.CreateOptions{
		CardID:   ip.CardID,
		Content:  ip.Content,
		ParentID: ip.ParentID,
	})

	if err != nil {
		uc.l.Errorf(ctx, "internal.comments.usecase.Create.repo.Create: %v", err)
		return comments.DetailOutput{}, err
	}

	u, err := uc.userUC.Detail(ctx, sc, c.UserID)
	if err != nil {
		if err == user.ErrUserNotFound {
			uc.l.Warnf(ctx, "internal.comments.usecase.Create.userUC.Detail: %v", err)
			return comments.DetailOutput{}, err
		}
		uc.l.Errorf(ctx, "internal.comments.usecase.Create.userUC.Detail: %v", err)
		return comments.DetailOutput{}, err
	}

	// Broadcast comment created event
	uc.broadcastCommentEvent(ctx, c.CardID, "comment_created", c, sc.UserID)

	return comments.DetailOutput{
		Comment: c,
		User:    u.User,
	}, nil
}

func (uc implUsecase) Update(ctx context.Context, sc models.Scope, ip comments.UpdateInput) (comments.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.comments.usecase.Update.repo.Detail.NotFound: %v", err)
			return comments.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.comments.usecase.Update.repo.Detail: %v", err)
		return comments.DetailOutput{}, err
	}

	c, err := uc.repo.Update(ctx, sc, repository.UpdateOptions{
		ID:       ip.ID,
		Content:  ip.Content,
		OldModel: oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.comments.usecase.Update.repo.Update: %v", err)
		return comments.DetailOutput{}, err
	}

	// Broadcast comment updated event
	uc.broadcastCommentEvent(ctx, c.CardID, "comment_updated", c, sc.UserID)

	return comments.DetailOutput{
		Comment: c,
	}, nil
}

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, ID string) (comments.DetailOutput, error) {
	c, err := uc.repo.Detail(ctx, sc, ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.comments.usecase.Detail.repo.Detail.NotFound: %v", err)
			return comments.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.comments.usecase.Detail.repo.Detail: %v", err)
		return comments.DetailOutput{}, err
	}
	return comments.DetailOutput{
		Comment: c,
	}, nil
}

func (uc implUsecase) Delete(ctx context.Context, sc models.Scope, ids []string) error {
	if len(ids) == 0 {
		uc.l.Warnf(ctx, "internal.comments.usecase.Delete.ids.Empty")
		return comments.ErrFieldRequired
	}

	err := uc.repo.Delete(ctx, sc, ids)
	if err != nil {
		uc.l.Errorf(ctx, "internal.comments.usecase.Delete.repo.Delete: %v", err)
		return err
	}
	return nil
}

func (uc implUsecase) GetByCard(ctx context.Context, sc models.Scope, cardID string) (comments.GetOutput, error) {
	// Verify card exists
	_, err := uc.cardsUC.Detail(ctx, sc, cardID)
	if err != nil {
		if err == cards.ErrCardNotFound {
			uc.l.Warnf(ctx, "internal.comments.usecase.GetByCard.cardsUC.Detail.CardNotFound: %v", err)
			return comments.GetOutput{}, comments.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.comments.usecase.GetByCard.cardsUC.Detail: %v", err)
		return comments.GetOutput{}, err
	}

	c, err := uc.repo.GetByCard(ctx, sc, cardID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.comments.usecase.GetByCard.repo.GetByCard: %v", err)
		return comments.GetOutput{}, err
	}

	uIDs := make([]string, len(c))
	for i, comment := range c {
		uIDs[i] = comment.UserID
	}
	uIDs = util.RemoveDuplicates(uIDs)
	us, err := uc.userUC.List(ctx, sc, user.ListInput{
		Filter: user.Filter{
			IDs: uIDs,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.comments.usecase.GetByCard.userUC.List: %v", err)
		return comments.GetOutput{}, err
	}

	return comments.GetOutput{
		Comments: c,
		Users:    us,
	}, nil
}
