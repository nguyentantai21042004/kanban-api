package usecase

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/boards"
	"gitlab.com/tantai-kanban/kanban-api/internal/boards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists"
	"gitlab.com/tantai-kanban/kanban-api/internal/role"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
	"gitlab.com/tantai-kanban/kanban-api/internal/websocket/service"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

type implUsecase struct {
	l      log.Logger
	repo   repository.Repository
	wsHub  *service.Hub
	userUC user.UseCase
	roleUC role.UseCase
	listUC lists.UseCase
	clock  func() time.Time
}

var _ boards.UseCase = &implUsecase{}

func New(l log.Logger, repo repository.Repository, wsHub *service.Hub, userUC user.UseCase, roleUC role.UseCase, listUC lists.UseCase) boards.UseCase {
	return &implUsecase{
		l:      l,
		repo:   repo,
		userUC: userUC,
		roleUC: roleUC,
		listUC: listUC,
		wsHub:  wsHub,
		clock:  util.Now,
	}
}

func (uc *implUsecase) SetList(listUC lists.UseCase) {
	uc.listUC = listUC
}
