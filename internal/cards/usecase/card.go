package usecase

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/cards"
	"gitlab.com/tantai-kanban/kanban-api/internal/cards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

// broadcastCardEvent broadcasts card events to WebSocket clients
func (uc implUsecase) broadcastCardEvent(ctx context.Context, boardID, eventType string, data interface{}, userID string) {
	if uc.wsHub == nil {
		return
	}

	uc.wsHub.BroadcastToBoard(boardID, eventType, data, userID)
}

func (uc implUsecase) Get(ctx context.Context, sc models.Scope, ip cards.GetInput) (cards.GetOutput, error) {
	u, p, err := uc.repo.Get(ctx, sc, repository.GetOptions{
		Filter:   ip.Filter,
		PagQuery: ip.PagQuery,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Get.repo.Get: %v", err)
		return cards.GetOutput{}, err
	}

	return cards.GetOutput{
		Cards:      u,
		Pagination: p,
	}, nil
}

func (uc implUsecase) Create(ctx context.Context, sc models.Scope, ip cards.CreateInput) (cards.DetailOutput, error) {
	// Get next position in list
	maxPosition, err := uc.repo.GetMaxPosition(ctx, sc, ip.ListID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Create.repo.GetMaxPosition: %v", err)
		return cards.DetailOutput{}, err
	}

	// Set default priority if not provided
	if ip.Priority == "" {
		ip.Priority = models.CardPriorityMedium
	}

	b, err := uc.repo.Create(ctx, sc, repository.CreateOptions{
		ListID:      ip.ListID,
		Title:       ip.Title,
		Description: ip.Description,
		Position:    maxPosition + 1.0,
		Priority:    ip.Priority,
		Labels:      ip.Labels,
		DueDate:     ip.DueDate,
	})

	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Create.repo.Create: %v", err)
		return cards.DetailOutput{}, err
	}

	// Get BoardID from ListID for broadcasting
	boardID, err := uc.repo.GetBoardIDFromListID(ctx, ip.ListID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Create.repo.GetBoardIDFromListID: %v", err)
		// Continue without broadcasting rather than failing the entire operation
	} else {
		// Broadcast card created event to board
		uc.broadcastCardEvent(ctx, boardID, "card_created", b, sc.UserID)
	}

	return cards.DetailOutput{
		Card: b,
	}, nil
}

func (uc implUsecase) Update(ctx context.Context, sc models.Scope, ip cards.UpdateInput) (cards.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.Update.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.Update.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

	b, err := uc.repo.Update(ctx, sc, repository.UpdateOptions{
		ID:          ip.ID,
		Title:       ip.Title,
		Description: ip.Description,
		Priority:    ip.Priority,
		Labels:      ip.Labels,
		DueDate:     ip.DueDate,
		OldModel:    oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Update.repo.Update: %v", err)
		return cards.DetailOutput{}, err
	}

	// Get BoardID from ListID for broadcasting
	boardID, err := uc.repo.GetBoardIDFromListID(ctx, b.ListID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Update.repo.GetBoardIDFromListID: %v", err)
		// Continue without broadcasting rather than failing the entire operation
	} else {
		// Broadcast card updated event to board
		uc.broadcastCardEvent(ctx, boardID, "card_updated", b, sc.UserID)
	}

	return cards.DetailOutput{
		Card: b,
	}, nil
}

func (uc implUsecase) Move(ctx context.Context, sc models.Scope, ip cards.MoveInput) (cards.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.Move.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.Move.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

	uc.l.Infof(ctx, "Moving card %s from list %s to list %s at position %d", ip.ID, oldModel.ListID, ip.ListID, ip.Position)

	_, err = uc.repo.Move(ctx, sc, repository.MoveOptions{
		ID:       ip.ID,
		ListID:   ip.ListID,
		Position: ip.Position,
		OldModel: oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Move.repo.Move: %v", err)
		return cards.DetailOutput{}, err
	}

	// Fetch the updated card to ensure we have complete data
	updatedCard, err := uc.repo.Detail(ctx, sc, ip.ID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Move.repo.Detail.AfterMove: %v", err)
		return cards.DetailOutput{}, err
	}

	uc.l.Infof(ctx, "Card moved successfully: %s, Title: %s, ListID: %s, Position: %f",
		updatedCard.ID, updatedCard.Title, updatedCard.ListID, updatedCard.Position)

	// Get BoardID from ListID for broadcasting
	boardID, err := uc.repo.GetBoardIDFromListID(ctx, updatedCard.ListID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Move.repo.GetBoardIDFromListID: %v", err)
		// Continue without broadcasting rather than failing the entire operation
	} else {
		// Broadcast card moved event to board
		uc.broadcastCardEvent(ctx, boardID, "card_moved", updatedCard, sc.UserID)
	}

	return cards.DetailOutput{
		Card: updatedCard,
	}, nil
}

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, ID string) (cards.DetailOutput, error) {
	b, err := uc.repo.Detail(ctx, sc, ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.Detail.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.Detail.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}
	return cards.DetailOutput{
		Card: b,
	}, nil
}

func (uc implUsecase) Delete(ctx context.Context, sc models.Scope, ids []string) error {
	if len(ids) == 0 {
		uc.l.Warnf(ctx, "internal.cards.usecase.Delete.ids.Empty")
		return cards.ErrFieldRequired
	}

	// Get cards before deletion for broadcasting
	cardsToDelete := make([]models.Card, 0, len(ids))
	for _, id := range ids {
		card, err := uc.repo.Detail(ctx, sc, id)
		if err != nil {
			uc.l.Warnf(ctx, "internal.cards.usecase.Delete.repo.Detail: %v", err)
			continue
		}
		cardsToDelete = append(cardsToDelete, card)
	}

	err := uc.repo.Delete(ctx, sc, ids)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Delete.repo.Delete: %v", err)
		return err
	}

	// Broadcast card deleted events
	for _, card := range cardsToDelete {
		// Get BoardID from ListID for broadcasting
		boardID, err := uc.repo.GetBoardIDFromListID(ctx, card.ListID)
		if err != nil {
			uc.l.Errorf(ctx, "internal.cards.usecase.Delete.repo.GetBoardIDFromListID: %v", err)
			continue
		}

		uc.broadcastCardEvent(ctx, boardID, "card_deleted", map[string]interface{}{
			"id": card.ID,
		}, sc.UserID)
	}

	return nil
}

func (uc implUsecase) GetActivities(ctx context.Context, sc models.Scope, ip cards.GetActivitiesInput) (cards.GetActivitiesOutput, error) {
	activities, err := uc.repo.GetActivities(ctx, sc, repository.GetActivitiesOptions{
		CardID: ip.CardID,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.GetActivities.repo.GetActivities: %v", err)
		return cards.GetActivitiesOutput{}, err
	}

	return cards.GetActivitiesOutput{
		Activities: activities,
	}, nil
}
