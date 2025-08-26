package usecase

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/boards"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/websocket"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

func (uc implUsecase) broadcastListEvent(ctx context.Context, boardID, eventType string, data interface{}, userID string) error {
	if uc.wsHub == nil {
		uc.l.Warnf(ctx, "internal.lists.usecase.broadcastListEvent.wsHub.BroadcastToBoard: wsHub is nil")
		return nil
	}

	err := uc.wsHub.BroadcastToBoard(ctx, boardID, eventType, data, userID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.broadcastListEvent.wsHub.BroadcastToBoard: %v", err)
		return err
	}

	return nil
}

func (uc implUsecase) Get(ctx context.Context, sc models.Scope, ip lists.GetInput) (lists.GetOutput, error) {
	// Only admin can see all lists
	me, err := uc.userUC.DetailMe(ctx, sc)
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Get.userUC.DetailMe: %v", err)
		return lists.GetOutput{}, err
	}

	rl, err := uc.roleUC.Detail(ctx, sc, me.User.RoleID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Get.roleUC.Detail: %v", err)
		return lists.GetOutput{}, err
	}

	if rl.Code != models.ADMIN_ROLE {
		ip.Filter.CreatedBy = me.User.ID
	}

	lsts, p, err := uc.repo.Get(ctx, sc, repository.GetOptions{
		Filter:   ip.Filter,
		PagQuery: ip.PagQuery,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Get.repo.Get: %v", err)
		return lists.GetOutput{}, err
	}

	return lists.GetOutput{
		Lists:      lsts,
		Pagination: p,
	}, nil
}

func (uc implUsecase) Create(ctx context.Context, sc models.Scope, ip lists.CreateInput) (lists.DetailOutput, error) {
	b, err := uc.boardUC.Detail(ctx, sc, ip.BoardID)
	if err != nil {
		if err == boards.ErrNotFound {
			uc.l.Warnf(ctx, "internal.lists.usecase.Create.boardUC.Detail.NotFound: %v", err)
			return lists.DetailOutput{}, boards.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.lists.usecase.Create.boardUC.Detail: %v", err)
		return lists.DetailOutput{}, err
	}

	pos, err := uc.repo.GetPosition(ctx, sc, repository.GetPositionOptions{
		BoardID: ip.BoardID,
		ASC:     false,
	})
	if err != nil {
		if err == repository.ErrNotFound {
			pos = ""
		} else {
			uc.l.Errorf(ctx, "internal.lists.usecase.Create.repo.GetPosition: %v", err)
			return lists.DetailOutput{}, err
		}
	}

	nwPst, err := uc.positionUC.GeneratePosition(pos, "")
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Create.positionUC.GeneratePosition: %v", err)
		return lists.DetailOutput{}, err
	}

	l, err := uc.repo.Create(ctx, sc, repository.CreateOptions{
		BoardID:  ip.BoardID,
		Name:     ip.Name,
		Position: nwPst,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Create.repo.Create: %v", err)
		return lists.DetailOutput{}, err
	}

	// Broadcast list created event
	err = uc.broadcastListEvent(ctx, ip.BoardID, websocket.MSG_LIST_CREATED, l, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Create.broadcastListEvent: %v", err)
	}

	return lists.DetailOutput{
		Board: b.Board,
		List:  l,
	}, nil
}

func (uc implUsecase) Update(ctx context.Context, sc models.Scope, ip lists.UpdateInput) (lists.DetailOutput, error) {
	om, err := uc.repo.Detail(ctx, sc, ip.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.lists.usecase.Update.repo.Detail.NotFound: %v", err)
			return lists.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.lists.usecase.Update.repo.Detail: %v", err)
		return lists.DetailOutput{}, err
	}

	_, err = uc.repo.Update(ctx, sc, repository.UpdateOptions{
		ID:       ip.ID,
		Name:     ip.Name,
		OldModel: om,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Update.repo.Update: %v", err)
		return lists.DetailOutput{}, err
	}

	l, err := uc.repo.Detail(ctx, sc, ip.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.lists.usecase.Update.repo.Detail.NotFound: %v", err)
			return lists.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.lists.usecase.Update.repo.Detail: %v", err)
		return lists.DetailOutput{}, err
	}

	err = uc.broadcastListEvent(ctx, l.BoardID, websocket.MSG_LIST_UPDATED, l, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Update.broadcastListEvent: %v", err)
	}

	return lists.DetailOutput{
		List: l,
	}, nil
}

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, ID string) (lists.DetailOutput, error) {
	b, err := uc.repo.Detail(ctx, sc, ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.lists.usecase.Detail.repo.Detail.NotFound: %v", err)
			return lists.DetailOutput{}, lists.ErrNotFound
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

	ls, err := uc.repo.List(ctx, sc, repository.ListOptions{
		Filter: lists.Filter{
			IDs: ids,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Delete.repo.Delete: %v", err)
		return err
	}

	if len(ls) != len(ids) {
		uc.l.Warnf(ctx, "internal.lists.usecase.Delete.repo.List.LengthMismatch: %v", err)
		return lists.ErrNotFound
	}

	err = uc.repo.Delete(ctx, sc, ids)
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Delete.repo.Delete: %v", err)
		return err
	}

	for _, list := range ls {
		err = uc.broadcastListEvent(ctx, list.BoardID, websocket.MSG_LIST_DELETED, map[string]interface{}{
			"id": list.ID,
		}, sc.UserID)
		if err != nil {
			uc.l.Errorf(ctx, "internal.lists.usecase.Delete.broadcastListEvent: %v", err)
		}
	}

	return nil
}

func (uc implUsecase) Move(ctx context.Context, sc models.Scope, ip lists.MoveInput) error {
	// Collect involved list IDs and de-duplicate
	listIDs := []string{ip.ID}
	if ip.AfterID != "" {
		listIDs = append(listIDs, ip.AfterID)
	}
	if ip.BeforeID != "" {
		listIDs = append(listIDs, ip.BeforeID)
	}
	listIDs = util.RemoveDuplicates(listIDs)

	// Fetch involved lists in one query
	ls, err := uc.repo.List(ctx, sc, repository.ListOptions{
		Filter: lists.Filter{IDs: listIDs},
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Move.repo.List: %v", err)
		return err
	}

	// Build map for quick lookup
	listMap := make(map[string]models.List, len(ls))
	for _, l := range ls {
		listMap[l.ID] = l
	}

	// Ensure target list exists
	if _, ok := listMap[ip.ID]; !ok {
		uc.l.Errorf(ctx, "internal.lists.usecase.Move.repo.List.NotFound: %v", repository.ErrNotFound)
		return repository.ErrNotFound
	}

	// Determine neighbor positions
	afterPos := ""
	beforePos := ""
	if after, ok := listMap[ip.AfterID]; ok {
		afterPos = after.Position
	}
	if before, ok := listMap[ip.BeforeID]; ok {
		beforePos = before.Position
	}

	// Generate new position
	newPos, err := uc.positionUC.GeneratePosition(afterPos, beforePos)
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Move.positionUC.GeneratePosition: %v", err)
		return err
	}

	// Apply move
	if _, err = uc.repo.Move(ctx, sc, repository.MoveOptions{
		ID:          ip.ID,
		BoardID:     ip.BoardID,
		NewPosition: newPos,
	}); err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Move.repo.Move: %v", err)
		return err
	}

	// Fetch updated list for broadcast consistency
	updList, err := uc.repo.Detail(ctx, sc, ip.ID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Move.repo.Detail.AfterMove: %v", err)
		return err
	}

	// Broadcast event
	if err := uc.broadcastListEvent(ctx, updList.BoardID, websocket.MSG_LIST_MOVED, updList, sc.UserID); err != nil {
		uc.l.Errorf(ctx, "internal.lists.usecase.Move.broadcastListEvent: %v", err)
	}
	return nil
}
