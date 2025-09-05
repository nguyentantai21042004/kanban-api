package usecase

import (
	"context"
	"sync"

	"github.com/nguyentantai21042004/kanban-api/internal/boards"
	"github.com/nguyentantai21042004/kanban-api/internal/cards"
	"github.com/nguyentantai21042004/kanban-api/internal/cards/repository"
	"github.com/nguyentantai21042004/kanban-api/internal/lists"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/internal/user"
	"github.com/nguyentantai21042004/kanban-api/internal/websocket"
	"github.com/nguyentantai21042004/kanban-api/pkg/util"
)

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, ID string) (cards.DetailOutput, error) {
	c, err := uc.repo.Detail(ctx, sc, ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.Detail.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.Detail.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

	var (
		ob boards.DetailOutput
		ol lists.DetailOutput
		// usrs    []models.User
		errChan = make(chan error, 2)
		wg      sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		ob, err = uc.boardUC.Detail(ctx, sc, c.BoardID)
		if err != nil {
			if err == repository.ErrNotFound {
				uc.l.Warnf(ctx, "internal.cards.usecase.Detail.boardUC.Detail.NotFound: %v", err)
				errChan <- repository.ErrNotFound
			}
			uc.l.Errorf(ctx, "internal.cards.usecase.Detail.boardUC.Detail: %v", err)
			errChan <- err
		}
		errChan <- nil
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		ol, err = uc.listUC.Detail(ctx, sc, c.ListID)
		if err != nil {
			if err == lists.ErrNotFound {
				uc.l.Warnf(ctx, "internal.cards.usecase.Detail.listUC.Detail.NotFound: %v", err)
				errChan <- lists.ErrNotFound
			}
			uc.l.Errorf(ctx, "internal.cards.usecase.Detail.listUC.Detail: %v", err)
			errChan <- err
		}
		errChan <- nil
	}()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	var err error
	// 	usrs, err = uc.userUC.List(ctx, sc, user.ListInput{
	// 		Filter: user.Filter{
	// 			IDs: []string{c.CreatedUserID, c.AssignedTo, c.UpdatedUserID},
	// 		},
	// 	})
	// 	if err != nil {
	// 		uc.l.Errorf(ctx, "internal.cards.usecase.Detail.userUC.List: %v", err)
	// 		errChan <- err
	// 	}
	// 	errChan <- nil
	// }()

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			uc.l.Errorf(ctx, "internal.cards.usecase.Detail.errChan: %v", err)
			return cards.DetailOutput{}, err
		}
	}

	return cards.DetailOutput{
		Card:  c,
		List:  ol.List,
		Board: ob.Board,
		// Users: usrs,
	}, nil
}

func (uc implUsecase) Get(ctx context.Context, sc models.Scope, ip cards.GetInput) (cards.GetOutput, error) {
	// Only admin can see all cards
	me, err := uc.userUC.DetailMe(ctx, sc)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Get.userUC.DetailMe: %v", err)
		return cards.GetOutput{}, err
	}

	rl, err := uc.roleUC.Detail(ctx, sc, me.User.RoleID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Get.roleUC.Detail: %v", err)
		return cards.GetOutput{}, err
	}

	if rl.Code != models.ADMIN_ROLE {
		ip.Filter.CreatedBy = me.User.ID
	}

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

func (uc implUsecase) Move(ctx context.Context, sc models.Scope, ip cards.MoveInput) error {
	cIDs := []string{ip.ID}
	if ip.AfterID != "" {
		cIDs = append(cIDs, ip.AfterID)
	}
	if ip.BeforeID != "" {
		cIDs = append(cIDs, ip.BeforeID)
	}
	cIDs = util.RemoveDuplicates(cIDs)

	cs, err := uc.repo.List(ctx, sc, repository.ListOptions{
		Filter: cards.Filter{
			IDs: cIDs,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Move.repo.List: %v", err)
		return err
	}

	// Build a map for quick lookup
	crdMap := make(map[string]models.Card, len(cs))
	for _, crd := range cs {
		crdMap[crd.ID] = crd
	}

	// Ensure the card to move exists
	if _, ok := crdMap[ip.ID]; !ok {
		uc.l.Errorf(ctx, "internal.cards.usecase.Move.repo.List.NotFound: %v", repository.ErrNotFound)
		return repository.ErrNotFound
	}

	// Get positions of after/before cards if they exist
	var afterPst, beforePst string = "", ""
	if afTrCrd, ok := crdMap[ip.AfterID]; ok {
		afterPst = afTrCrd.Position
	}
	if bfTrCrd, ok := crdMap[ip.BeforeID]; ok {
		beforePst = bfTrCrd.Position
	}

	// Generate new position for the card
	nwPst, err := uc.positionUC.GeneratePosition(afterPst, beforePst)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Move.positionUC.GeneratePosition: %v", err)
		return err
	}

	// Move the card in the repository
	if _, err = uc.repo.Move(ctx, sc, repository.MoveOptions{
		ID:          ip.ID,
		ListID:      ip.ListID,
		NewPosition: nwPst,
	}); err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Move.repo.Move: %v", err)
		return err
	}

	// Fetch the updated card to ensure we have complete data
	updCard, err := uc.repo.Detail(ctx, sc, ip.ID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Move.repo.Detail.AfterMove: %v", err)
		return err
	}

	// Broadcast the card move event, but don't fail the operation if broadcasting fails
	if err := uc.wsHub.BroadcastToBoard(ctx, updCard.BoardID, websocket.MSG_CARD_MOVED, updCard, sc.UserID); err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Move.wsHub.BroadcastToBoard: %v", err)
	}

	return nil
}

func (uc implUsecase) Create(ctx context.Context, sc models.Scope, ip cards.CreateInput) (cards.DetailOutput, error) {
	var (
		ob      boards.DetailOutput
		ol      lists.DetailOutput
		errChan = make(chan error, 2)
		wg      sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		ob, err = uc.boardUC.Detail(ctx, sc, ip.BoardID)
		if err != nil {
			if err == repository.ErrNotFound {
				uc.l.Warnf(ctx, "internal.cards.usecase.Create.boardUC.Detail.NotFound: %v", err)
				errChan <- repository.ErrNotFound
			}
			uc.l.Errorf(ctx, "internal.cards.usecase.Create.boardUC.Detail: %v", err)
			errChan <- err
		}
		errChan <- nil
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		ol, err = uc.listUC.Detail(ctx, sc, ip.ListID)
		if err != nil {
			if err == lists.ErrNotFound {
				uc.l.Warnf(ctx, "internal.cards.usecase.Create.listUC.Detail.NotFound: %v", err)
				errChan <- lists.ErrNotFound
			}
			uc.l.Errorf(ctx, "internal.cards.usecase.Create.listUC.Detail: %v", err)
			errChan <- err
		}
		errChan <- nil
	}()

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			uc.l.Errorf(ctx, "internal.cards.usecase.Create.errChan: %v", err)
			return cards.DetailOutput{}, err
		}
	}

	// Get current max position in list
	mxPst, err := uc.repo.GetPosition(ctx, sc, repository.GetPositionOptions{
		ListID: ip.ListID,
		ASC:    false,
	})
	if err != nil && err != repository.ErrNotFound {
		uc.l.Errorf(ctx, "internal.cards.usecase.Create.repo.GetMaxPosition: %v", err)
		return cards.DetailOutput{}, err
	}

	pst, err := uc.positionUC.GeneratePosition(mxPst, "")
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Create.positionUC.GeneratePosition: %v", err)
		return cards.DetailOutput{}, err
	}

	// Set default priority if not provided
	if ip.Priority == "" {
		ip.Priority = models.CardPriorityMedium
	}

	checklist := make([]repository.CheckListOptions, len(ip.Checklist))
	for i, c := range ip.Checklist {
		checklist[i] = repository.CheckListOptions{
			Content:     c.Content,
			IsCompleted: c.IsCompleted,
		}
	}

	b, err := uc.repo.Create(ctx, sc, repository.CreateOptions{
		BoardID:        ip.BoardID,
		ListID:         ip.ListID,
		Name:           ip.Name,
		Alias:          util.BuildAlias(ip.Name),
		Description:    ip.Description,
		Position:       pst,
		Priority:       ip.Priority,
		Labels:         ip.Labels,
		DueDate:        ip.DueDate,
		CreatedBy:      sc.UserID,
		AssignedTo:     ip.AssignedTo,
		EstimatedHours: ip.EstimatedHours,
		StartDate:      ip.StartDate,
		Tags:           ip.Tags,
		Checklist:      checklist,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Create.repo.Create: %v", err)
		return cards.DetailOutput{}, err
	}

	err = uc.wsHub.BroadcastToBoard(ctx, ob.Board.ID, websocket.MSG_CARD_CREATED, b, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Create.wsHub.BroadcastToBoard: %v", err)
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
		Board: ob.Board,
		Users: usrs,
	}, nil
}

func (uc implUsecase) Update(ctx context.Context, sc models.Scope, ip cards.UpdateInput) (cards.DetailOutput, error) {
	// First get the existing card to fetch its board and list IDs
	oc, err := uc.repo.Detail(ctx, sc, ip.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.cards.usecase.Update.repo.Detail.NotFound: %v", err)
			return cards.DetailOutput{}, repository.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.cards.usecase.Update.repo.Detail: %v", err)
		return cards.DetailOutput{}, err
	}

	var (
		ob      boards.DetailOutput
		ol      lists.DetailOutput
		errChan = make(chan error, 2)
		wg      sync.WaitGroup
	)

	// Use the existing card's board and list IDs for validation
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		ob, err = uc.boardUC.Detail(ctx, sc, oc.BoardID)
		if err != nil {
			if err == repository.ErrNotFound {
				uc.l.Warnf(ctx, "internal.cards.usecase.Update.boardUC.Detail.NotFound: %v", err)
				errChan <- repository.ErrNotFound
			}
			uc.l.Errorf(ctx, "internal.cards.usecase.Update.boardUC.Detail: %v", err)
			errChan <- err
		}
		errChan <- nil
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		ol, err = uc.listUC.Detail(ctx, sc, oc.ListID)
		if err != nil {
			if err == lists.ErrNotFound {
				uc.l.Warnf(ctx, "internal.cards.usecase.Update.listUC.Detail.NotFound: %v", err)
				errChan <- lists.ErrNotFound
			}
			uc.l.Errorf(ctx, "internal.cards.usecase.Update.listUC.Detail: %v", err)
			errChan <- err
		}
		errChan <- nil
	}()

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			uc.l.Errorf(ctx, "internal.cards.usecase.Update.errChan: %v", err)
			return cards.DetailOutput{}, err
		}
	}

	b, err := uc.repo.Update(ctx, sc, repository.UpdateOptions{
		ID:             ip.ID,
		Name:           ip.Name,
		Alias:          util.BuildAlias(ip.Name),
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
		OldModel:       oc,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Update.repo.Update: %v", err)
		return cards.DetailOutput{}, err
	}

	err = uc.wsHub.BroadcastToBoard(ctx, oc.BoardID, websocket.MSG_CARD_UPDATED, b, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Update.wsHub.BroadcastToBoard: %v", err)
	}

	return cards.DetailOutput{
		Card:  b,
		List:  ol.List,
		Board: ob.Board,
	}, nil
}

func (uc implUsecase) Delete(ctx context.Context, sc models.Scope, ids []string) error {
	if len(ids) == 0 {
		uc.l.Errorf(ctx, "internal.cards.usecase.Delete.ids.Empty")
		return cards.ErrFieldRequired
	}

	cs, err := uc.repo.List(ctx, sc, repository.ListOptions{
		Filter: cards.Filter{
			IDs: ids,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Delete.repo.Detail: %v", err)
		return err
	}

	if len(cs) != len(ids) {
		uc.l.Errorf(ctx, "internal.cards.usecase.Delete.repo.List.Empty")
		return cards.ErrCardNotFound
	}

	err = uc.repo.Delete(ctx, sc, ids)
	if err != nil {
		uc.l.Errorf(ctx, "internal.cards.usecase.Delete.repo.Delete: %v", err)
		return err
	}

	for _, card := range cs {
		err = uc.wsHub.BroadcastToBoard(ctx, card.BoardID, websocket.MSG_CARD_DELETED, card, sc.UserID)
		if err != nil {
			uc.l.Errorf(ctx, "internal.cards.usecase.Delete.wsHub.BroadcastToBoard: %v", err)
		}
	}

	return nil
}
