package usecase

import (
	"context"

	"github.com/nguyentantai21042004/kanban-api/internal/boards"
	"github.com/nguyentantai21042004/kanban-api/internal/boards/repository"
	"github.com/nguyentantai21042004/kanban-api/internal/lists"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/internal/user"
	"github.com/nguyentantai21042004/kanban-api/pkg/util"
)

func (uc implUsecase) broadcastBoardEvent(ctx context.Context, boardID, eventType string, data interface{}, userID string) error {
	if uc.wsHub == nil {
		uc.l.Warnf(ctx, "internal.boards.usecase.broadcastBoardEvent.wsHub.BroadcastToBoard: wsHub is nil")
		return nil
	}

	err := uc.wsHub.BroadcastToBoard(ctx, boardID, eventType, data, userID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.boards.usecase.broadcastBoardEvent.wsHub.BroadcastToBoard: %v", err)
		return err
	}

	return nil
}

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

	// Create default lists for the new board
	defaultLists := []string{
		"Backlog",          // 1. Backlog
		"Data Preparation", // 2. Data Preparation
		"In Progress",      // 3. In Progress
		"Blocked",          // 4. Blocked
		"Completed",        // 5. Completed
	}

	// Create default lists sequentially to ensure proper positioning
	for i, listName := range defaultLists {
		_, err := uc.listUC.Create(ctx, sc, lists.CreateInput{
			BoardID: b.ID,
			Name:    listName,
		})
		if err != nil {
			uc.l.Errorf(ctx, "internal.boards.usecase.Create.listUC.Create.defaultList[%d]: %v", i, err)
			// Don't fail the board creation if list creation fails
			// Just log the error and continue
		}
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
		Users: []models.User{u.User},
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

	// Broadcast board updated event
	err = uc.broadcastBoardEvent(ctx, b.ID, "board_updated", b, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.boards.usecase.Update.broadcastBoardEvent: %v", err)
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
			return boards.DetailOutput{}, boards.ErrNotFound
		}
		uc.l.Errorf(ctx, "internal.boards.usecase.Detail.repo.Detail: %v", err)
		return boards.DetailOutput{}, err
	}

	uIDs := []string{*b.CreatedBy}
	us, err := uc.userUC.List(ctx, sc, user.ListInput{
		Filter: user.Filter{
			IDs: uIDs,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.boards.usecase.Detail.userUC.List: %v", err)
		return boards.DetailOutput{}, err
	}

	return boards.DetailOutput{
		Board: b,
		Users: us,
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

// Dashboard: total via repo.Get; active computed in admin via cards dashboard
func (uc implUsecase) Dashboard(ctx context.Context, sc models.Scope, ip boards.DashboardInput) (boards.BoardsDashboardOutput, error) {
	bs, _, err := uc.repo.List(ctx, sc, repository.ListOptions{})
	if err != nil {
		uc.l.Errorf(ctx, "internal.boards.usecase.Dashboard.repo.Get: %v", err)
		return boards.BoardsDashboardOutput{}, err
	}
	return boards.BoardsDashboardOutput{Total: int64(len(bs)), Active: 0}, nil
}
