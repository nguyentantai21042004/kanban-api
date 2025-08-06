package usecase

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/cards"
	"gitlab.com/tantai-kanban/kanban-api/internal/cards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
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
	// Validate List
	ol, err := uc.listUC.Detail(ctx, sc, ip.ListID)
	if err != nil {
		if err == lists.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.Create.listUC.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, lists.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.Create.listUC.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

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
		ListID:         ip.ListID,
		Name:           ip.Name,
		Description:    ip.Description,
		Position:       maxPosition + 1.0,
		Priority:       ip.Priority,
		Labels:         ip.Labels,
		DueDate:        ip.DueDate,
		CreatedBy:      sc.UserID,
		AssignedTo:     ip.AssignedTo,
		EstimatedHours: ip.EstimatedHours,
		StartDate:      ip.StartDate,
		Tags:           ip.Tags,
		Checklist:      ip.Checklist,
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

	userIDs := []string{sc.UserID}
	usrs, err := uc.userUC.List(ctx, sc, user.ListInput{
		Filter: user.Filter{
			IDs: userIDs,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Create.userUC.Get: %v", err)
		return cards.DetailOutput{}, err
	}

	return cards.DetailOutput{
		Card:  b,
		List:  ol.List,
		Users: usrs,
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
		ID:             ip.ID,
		Name:           ip.Name,
		Description:    ip.Description,
		Priority:       ip.Priority,
		Labels:         ip.Labels,
		DueDate:        ip.DueDate,
		AssignedTo:     ip.AssignedTo,
		EstimatedHours: ip.EstimatedHours,
		ActualHours:    ip.ActualHours,
		StartDate:      ip.StartDate,
		CompletionDate: ip.CompletionDate,
		Tags:           ip.Tags,
		Checklist:      ip.Checklist,
		OldModel:       oldModel,
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

	uc.l.Infof(ctx, "Card moved successfully: %s, Name: %s, ListID: %s, Position: %f",
		updatedCard.ID, updatedCard.Name, updatedCard.ListID, updatedCard.Position)

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

// New methods for enhanced functionality
func (uc implUsecase) Assign(ctx context.Context, sc models.Scope, ip cards.AssignInput) (cards.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.Assign.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.Assign.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

	card, err := uc.repo.Assign(ctx, sc, repository.AssignOptions{
		CardID:     ip.CardID,
		AssignedTo: ip.AssignedTo,
		OldModel:   oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Assign.repo.Assign: %v", err)
		return cards.DetailOutput{}, err
	}

	// Broadcast card assigned event
	boardID, err := uc.repo.GetBoardIDFromListID(ctx, card.ListID)
	if err == nil {
		uc.broadcastCardEvent(ctx, boardID, "card_assigned", card, sc.UserID)
	}

	return cards.DetailOutput{
		Card: card,
	}, nil
}

func (uc implUsecase) Unassign(ctx context.Context, sc models.Scope, ip cards.UnassignInput) (cards.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.Unassign.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.Unassign.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

	card, err := uc.repo.Unassign(ctx, sc, repository.UnassignOptions{
		CardID:   ip.CardID,
		OldModel: oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Unassign.repo.Unassign: %v", err)
		return cards.DetailOutput{}, err
	}

	// Broadcast card unassigned event
	boardID, err := uc.repo.GetBoardIDFromListID(ctx, card.ListID)
	if err == nil {
		uc.broadcastCardEvent(ctx, boardID, "card_unassigned", card, sc.UserID)
	}

	return cards.DetailOutput{
		Card: card,
	}, nil
}

func (uc implUsecase) AddAttachment(ctx context.Context, sc models.Scope, ip cards.AddAttachmentInput) (cards.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.AddAttachment.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.AddAttachment.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

	card, err := uc.repo.AddAttachment(ctx, sc, repository.AddAttachmentOptions{
		CardID:       ip.CardID,
		AttachmentID: ip.AttachmentID,
		OldModel:     oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.AddAttachment.repo.AddAttachment: %v", err)
		return cards.DetailOutput{}, err
	}

	// Broadcast card attachment added event
	boardID, err := uc.repo.GetBoardIDFromListID(ctx, card.ListID)
	if err == nil {
		uc.broadcastCardEvent(ctx, boardID, "card_attachment_added", card, sc.UserID)
	}

	return cards.DetailOutput{
		Card: card,
	}, nil
}

func (uc implUsecase) RemoveAttachment(ctx context.Context, sc models.Scope, ip cards.RemoveAttachmentInput) (cards.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.RemoveAttachment.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.RemoveAttachment.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

	card, err := uc.repo.RemoveAttachment(ctx, sc, repository.RemoveAttachmentOptions{
		CardID:       ip.CardID,
		AttachmentID: ip.AttachmentID,
		OldModel:     oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.RemoveAttachment.repo.RemoveAttachment: %v", err)
		return cards.DetailOutput{}, err
	}

	// Broadcast card attachment removed event
	boardID, err := uc.repo.GetBoardIDFromListID(ctx, card.ListID)
	if err == nil {
		uc.broadcastCardEvent(ctx, boardID, "card_attachment_removed", card, sc.UserID)
	}

	return cards.DetailOutput{
		Card: card,
	}, nil
}

func (uc implUsecase) UpdateTimeTracking(ctx context.Context, sc models.Scope, ip cards.UpdateTimeTrackingInput) (cards.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.UpdateTimeTracking.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.UpdateTimeTracking.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

	card, err := uc.repo.UpdateTimeTracking(ctx, sc, repository.UpdateTimeTrackingOptions{
		CardID:         ip.CardID,
		EstimatedHours: ip.EstimatedHours,
		ActualHours:    ip.ActualHours,
		OldModel:       oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.UpdateTimeTracking.repo.UpdateTimeTracking: %v", err)
		return cards.DetailOutput{}, err
	}

	// Broadcast card time tracking updated event
	boardID, err := uc.repo.GetBoardIDFromListID(ctx, card.ListID)
	if err == nil {
		uc.broadcastCardEvent(ctx, boardID, "card_time_tracking_updated", card, sc.UserID)
	}

	return cards.DetailOutput{
		Card: card,
	}, nil
}

func (uc implUsecase) UpdateChecklist(ctx context.Context, sc models.Scope, ip cards.UpdateChecklistInput) (cards.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.UpdateChecklist.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.UpdateChecklist.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

	card, err := uc.repo.UpdateChecklist(ctx, sc, repository.UpdateChecklistOptions{
		CardID:    ip.CardID,
		Checklist: ip.Checklist,
		OldModel:  oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.UpdateChecklist.repo.UpdateChecklist: %v", err)
		return cards.DetailOutput{}, err
	}

	// Broadcast card checklist updated event
	boardID, err := uc.repo.GetBoardIDFromListID(ctx, card.ListID)
	if err == nil {
		uc.broadcastCardEvent(ctx, boardID, "card_checklist_updated", card, sc.UserID)
	}

	return cards.DetailOutput{
		Card: card,
	}, nil
}

func (uc implUsecase) AddTag(ctx context.Context, sc models.Scope, ip cards.AddTagInput) (cards.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.AddTag.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.AddTag.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

	card, err := uc.repo.AddTag(ctx, sc, repository.AddTagOptions{
		CardID:   ip.CardID,
		Tag:      ip.Tag,
		OldModel: oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.AddTag.repo.AddTag: %v", err)
		return cards.DetailOutput{}, err
	}

	// Broadcast card tag added event
	boardID, err := uc.repo.GetBoardIDFromListID(ctx, card.ListID)
	if err == nil {
		uc.broadcastCardEvent(ctx, boardID, "card_tag_added", card, sc.UserID)
	}

	return cards.DetailOutput{
		Card: card,
	}, nil
}

func (uc implUsecase) RemoveTag(ctx context.Context, sc models.Scope, ip cards.RemoveTagInput) (cards.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.RemoveTag.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.RemoveTag.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

	card, err := uc.repo.RemoveTag(ctx, sc, repository.RemoveTagOptions{
		CardID:   ip.CardID,
		Tag:      ip.Tag,
		OldModel: oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.RemoveTag.repo.RemoveTag: %v", err)
		return cards.DetailOutput{}, err
	}

	// Broadcast card tag removed event
	boardID, err := uc.repo.GetBoardIDFromListID(ctx, card.ListID)
	if err == nil {
		uc.broadcastCardEvent(ctx, boardID, "card_tag_removed", card, sc.UserID)
	}

	return cards.DetailOutput{
		Card: card,
	}, nil
}

func (uc implUsecase) SetStartDate(ctx context.Context, sc models.Scope, ip cards.SetStartDateInput) (cards.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.SetStartDate.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.SetStartDate.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

	card, err := uc.repo.SetStartDate(ctx, sc, repository.SetStartDateOptions{
		CardID:    ip.CardID,
		StartDate: ip.StartDate,
		OldModel:  oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.SetStartDate.repo.SetStartDate: %v", err)
		return cards.DetailOutput{}, err
	}

	// Broadcast card start date set event
	boardID, err := uc.repo.GetBoardIDFromListID(ctx, card.ListID)
	if err == nil {
		uc.broadcastCardEvent(ctx, boardID, "card_start_date_set", card, sc.UserID)
	}

	return cards.DetailOutput{
		Card: card,
	}, nil
}

func (uc implUsecase) SetCompletionDate(ctx context.Context, sc models.Scope, ip cards.SetCompletionDateInput) (cards.DetailOutput, error) {
	oldModel, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.SetCompletionDate.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.SetCompletionDate.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

	card, err := uc.repo.SetCompletionDate(ctx, sc, repository.SetCompletionDateOptions{
		CardID:         ip.CardID,
		CompletionDate: ip.CompletionDate,
		OldModel:       oldModel,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.SetCompletionDate.repo.SetCompletionDate: %v", err)
		return cards.DetailOutput{}, err
	}

	// Broadcast card completion date set event
	boardID, err := uc.repo.GetBoardIDFromListID(ctx, card.ListID)
	if err == nil {
		uc.broadcastCardEvent(ctx, boardID, "card_completion_date_set", card, sc.UserID)
	}

	return cards.DetailOutput{
		Card: card,
	}, nil
}
