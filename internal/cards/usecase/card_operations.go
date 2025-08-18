package usecase

import (
	"context"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/cards"
	"gitlab.com/tantai-kanban/kanban-api/internal/cards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/websocket"
)

func (uc implUsecase) GetActivities(ctx context.Context, sc models.Scope, ip cards.GetActivitiesInput) (cards.GetActivitiesOutput, error) {
	c, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.GetActivities.repo.Detail.NotFound: %v", err)
			return cards.GetActivitiesOutput{}, cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.GetActivities.repo.Detail: %v", err)
		return cards.GetActivitiesOutput{}, err
	}

	atvs, pag, err := uc.repo.GetActivities(ctx, sc, repository.GetActivitiesOptions{
		Filter: repository.ActivityFilter{
			CardID: ip.CardID,
		},
		PagQuery: ip.PagQuery,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.GetActivities.repo.GetActivities: %v", err)
		return cards.GetActivitiesOutput{}, err
	}

	return cards.GetActivitiesOutput{
		Card:       c,
		Activities: atvs,
		Pagination: pag,
	}, nil
}

func (uc implUsecase) Assign(ctx context.Context, sc models.Scope, ip cards.AssignInput) error {
	om, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.Assign.repo.Detail.NotFound: %v", err)
			return cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.Assign.repo.Detail: %v", err)
		return err
	}

	usr, err := uc.userUC.Detail(ctx, sc, ip.AssignedTo)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Assign.userUC.List: %v", err)
		return err
	}

	_, err = uc.repo.Assign(ctx, sc, repository.AssignOptions{
		CardID:     ip.CardID,
		AssignedTo: usr.User.ID,
		OldModel:   om,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Assign.repo.Assign: %v", err)
		return err
	}

	crd, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Assign.repo.Detail: %v", err)
		return err
	}

	err = uc.wsHub.BroadcastToBoard(ctx, om.BoardID, websocket.MSG_CARD_UPDATED, crd, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Assign.wsHub.BroadcastToBoard: %v", err)
	}

	return nil
}

func (uc implUsecase) Unassign(ctx context.Context, sc models.Scope, ip cards.UnassignInput) error {
	om, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.Unassign.repo.Detail.NotFound: %v", err)
			return cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.Unassign.repo.Detail: %v", err)
		return err
	}

	// Perform unassign operation
	card, err := uc.repo.Unassign(ctx, sc, repository.UnassignOptions{
		CardID:   ip.CardID,
		OldModel: om,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Unassign.repo.Unassign: %v", err)
		return err
	}

	err = uc.wsHub.BroadcastToBoard(ctx, om.BoardID, websocket.MSG_CARD_UPDATED, card, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Unassign.wsHub.BroadcastToBoard: %v", err)
	}

	return nil
}

func (uc implUsecase) AddAttachment(ctx context.Context, sc models.Scope, ip cards.AddAttachmentInput) error {
	om, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.AddAttachment.repo.Detail.NotFound: %v", err)
			return cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.AddAttachment.repo.Detail: %v", err)
		return err
	}

	card, err := uc.repo.AddAttachment(ctx, sc, repository.AddAttachmentOptions{
		CardID:       ip.CardID,
		AttachmentID: ip.AttachmentID,
		OldModel:     om,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.AddAttachment.repo.AddAttachment: %v", err)
		return err
	}

	err = uc.wsHub.BroadcastToBoard(ctx, om.BoardID, websocket.MSG_CARD_UPDATED, card, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.AddAttachment.wsHub.BroadcastToBoard: %v", err)
	}

	return nil
}

func (uc implUsecase) RemoveAttachment(ctx context.Context, sc models.Scope, ip cards.RemoveAttachmentInput) error {
	om, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.RemoveAttachment.repo.Detail.NotFound: %v", err)
			return cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.RemoveAttachment.repo.Detail: %v", err)
		return err
	}

	card, err := uc.repo.RemoveAttachment(ctx, sc, repository.RemoveAttachmentOptions{
		CardID:       ip.CardID,
		AttachmentID: ip.AttachmentID,
		OldModel:     om,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.RemoveAttachment.repo.RemoveAttachment: %v", err)
		return err
	}

	err = uc.wsHub.BroadcastToBoard(ctx, om.BoardID, websocket.MSG_CARD_UPDATED, card, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.RemoveAttachment.wsHub.BroadcastToBoard: %v", err)
	}

	return nil
}

func (uc implUsecase) UpdateTimeTracking(ctx context.Context, sc models.Scope, ip cards.UpdateTimeTrackingInput) error {
	om, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.UpdateTimeTracking.repo.Detail.NotFound: %v", err)
			return cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.UpdateTimeTracking.repo.Detail: %v", err)
		return err
	}

	card, err := uc.repo.UpdateTimeTracking(ctx, sc, repository.UpdateTimeTrackingOptions{
		CardID:         ip.CardID,
		EstimatedHours: ip.EstimatedHours,
		ActualHours:    ip.ActualHours,
		OldModel:       om,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.UpdateTimeTracking.repo.UpdateTimeTracking: %v", err)
		return err
	}

	// Broadcast card time tracking updated event
	err = uc.wsHub.BroadcastToBoard(ctx, om.BoardID, websocket.MSG_CARD_UPDATED, card, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.UpdateTimeTracking.wsHub.BroadcastToBoard: %v", err)
	}

	return nil
}

func (uc implUsecase) UpdateChecklist(ctx context.Context, sc models.Scope, ip cards.UpdateChecklistInput) error {
	om, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.UpdateChecklist.repo.Detail.NotFound: %v", err)
			return cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.UpdateChecklist.repo.Detail: %v", err)
		return err
	}

	card, err := uc.repo.UpdateChecklist(ctx, sc, repository.UpdateChecklistOptions{
		CardID:    ip.CardID,
		Checklist: ip.Checklist,
		OldModel:  om,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.UpdateChecklist.repo.UpdateChecklist: %v", err)
		return err
	}

	// Broadcast card checklist updated event
	err = uc.wsHub.BroadcastToBoard(ctx, om.BoardID, websocket.MSG_CARD_UPDATED, card, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.UpdateChecklist.wsHub.BroadcastToBoard: %v", err)
	}

	return nil
}

func (uc implUsecase) AddTag(ctx context.Context, sc models.Scope, ip cards.AddTagInput) error {
	om, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.AddTag.repo.Detail.NotFound: %v", err)
			return cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.AddTag.repo.Detail: %v", err)
		return err
	}

	card, err := uc.repo.AddTag(ctx, sc, repository.AddTagOptions{
		CardID:   ip.CardID,
		Tag:      ip.Tag,
		OldModel: om,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.AddTag.repo.AddTag: %v", err)
		return err
	}

	err = uc.wsHub.BroadcastToBoard(ctx, om.BoardID, websocket.MSG_CARD_UPDATED, card, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.AddTag.wsHub.BroadcastToBoard: %v", err)
	}

	return nil
}

func (uc implUsecase) RemoveTag(ctx context.Context, sc models.Scope, ip cards.RemoveTagInput) error {
	om, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.RemoveTag.repo.Detail.NotFound: %v", err)
			return cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.RemoveTag.repo.Detail: %v", err)
		return err
	}

	card, err := uc.repo.RemoveTag(ctx, sc, repository.RemoveTagOptions{
		CardID:   ip.CardID,
		Tag:      ip.Tag,
		OldModel: om,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.RemoveTag.repo.RemoveTag: %v", err)
		return err
	}

	err = uc.wsHub.BroadcastToBoard(ctx, om.BoardID, websocket.MSG_CARD_UPDATED, card, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.RemoveTag.wsHub.BroadcastToBoard: %v", err)
	}

	return nil
}

func (uc implUsecase) SetStartDate(ctx context.Context, sc models.Scope, ip cards.SetStartDateInput) error {
	om, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.SetStartDate.repo.Detail.NotFound: %v", err)
			return cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.SetStartDate.repo.Detail: %v", err)
		return err
	}

	card, err := uc.repo.SetStartDate(ctx, sc, repository.SetStartDateOptions{
		CardID:    ip.CardID,
		StartDate: ip.StartDate,
		OldModel:  om,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.SetStartDate.repo.SetStartDate: %v", err)
		return err
	}

	err = uc.wsHub.BroadcastToBoard(ctx, om.BoardID, websocket.MSG_CARD_UPDATED, card, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.SetStartDate.wsHub.BroadcastToBoard: %v", err)
	}

	return nil
}

func (uc implUsecase) SetCompletionDate(ctx context.Context, sc models.Scope, ip cards.SetCompletionDateInput) error {
	om, err := uc.repo.Detail(ctx, sc, ip.CardID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.SetCompletionDate.repo.Detail.NotFound: %v", err)
			return cards.ErrCardNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.SetCompletionDate.repo.Detail: %v", err)
		return err
	}

	card, err := uc.repo.SetCompletionDate(ctx, sc, repository.SetCompletionDateOptions{
		CardID:         ip.CardID,
		CompletionDate: ip.CompletionDate,
		OldModel:       om,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.SetCompletionDate.repo.SetCompletionDate: %v", err)
		return err
	}

	err = uc.wsHub.BroadcastToBoard(ctx, om.BoardID, websocket.MSG_CARD_UPDATED, card, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.SetCompletionDate.wsHub.BroadcastToBoard: %v", err)
	}

	return nil
}

// Dashboard aggregates for cards module
func (uc implUsecase) Dashboard(ctx context.Context, sc models.Scope, ip cards.DashboardInput) (cards.CardsDashboardOutput, error) {
	// total
	totalList, err := uc.repo.List(ctx, sc, repository.ListOptions{Filter: cards.Filter{}})
	if err != nil {
		return cards.CardsDashboardOutput{}, err
	}

	// completed in range
	completedFrom := ip.From
	completedTo := ip.To
	completedList, err := uc.repo.List(ctx, sc, repository.ListOptions{Filter: cards.Filter{
		CompletionDateFrom: &completedFrom,
		CompletionDateTo:   &completedTo,
	}})
	if err != nil {
		return cards.CardsDashboardOutput{}, err
	}

	// overdue at ip.To
	toRef := ip.To
	overdueList, err := uc.repo.List(ctx, sc, repository.ListOptions{Filter: cards.Filter{
		UncompletedOnly: true,
		DueDateTo:       &toRef,
	}})
	if err != nil {
		return cards.CardsDashboardOutput{}, err
	}

	// fetch cards updated within period for active users/boards
	updatedFrom := ip.From
	updatedTo := ip.To
	updatedList, err := uc.repo.List(ctx, sc, repository.ListOptions{Filter: cards.Filter{UpdatedFrom: &updatedFrom, UpdatedTo: &updatedTo}})
	if err != nil {
		return cards.CardsDashboardOutput{}, err
	}
	activeUserSet := map[string]struct{}{}
	activeBoardSet := map[string]struct{}{}
	for _, c := range updatedList {
		if c.UpdatedBy != nil && *c.UpdatedBy != "" {
			activeUserSet[*c.UpdatedBy] = struct{}{}
		}
		if c.CreatedBy != nil && *c.CreatedBy != "" {
			activeUserSet[*c.CreatedBy] = struct{}{}
		}
		if c.BoardID != "" {
			activeBoardSet[c.BoardID] = struct{}{}
		}
	}
	activeUsers := make([]string, 0, len(activeUserSet))
	for id := range activeUserSet {
		activeUsers = append(activeUsers, id)
	}
	activeBoards := make([]string, 0, len(activeBoardSet))
	for id := range activeBoardSet {
		activeBoards = append(activeBoards, id)
	}

	// activity per day
	var activity []cards.ActivityPoint
	day := ip.From
	for !day.After(ip.To) {
		next := day.AddDate(0, 0, 1)
		// created
		from := day
		to := next.Add(-time.Nanosecond)
		createdList, err := uc.repo.List(ctx, sc, repository.ListOptions{Filter: cards.Filter{CreatedFrom: &from, CreatedTo: &to}})
		if err != nil {
			return cards.CardsDashboardOutput{}, err
		}
		// completed
		compList, err := uc.repo.List(ctx, sc, repository.ListOptions{Filter: cards.Filter{CompletionDateFrom: &from, CompletionDateTo: &to}})
		if err != nil {
			return cards.CardsDashboardOutput{}, err
		}
		activity = append(activity, cards.ActivityPoint{
			Date:           day.Format("2006-01-02"),
			CardsCreated:   int64(len(createdList)),
			CardsCompleted: int64(len(compList)),
		})
		day = next
	}

	return cards.CardsDashboardOutput{
		Total:          int64(len(totalList)),
		Completed:      int64(len(completedList)),
		Overdue:        int64(len(overdueList)),
		Activity:       activity,
		ActiveUserIDs:  activeUsers,
		ActiveBoardIDs: activeBoards,
	}, nil
}
