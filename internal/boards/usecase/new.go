package usecase

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/boards"
	"gitlab.com/tantai-kanban/kanban-api/internal/boards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/role"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
	"gitlab.com/tantai-kanban/kanban-api/internal/websocket/service"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

type implUsecase struct {
	l      log.Logger
	clock  func() time.Time
	repo   repository.Repository
	userUC user.UseCase
	roleUC role.UseCase
	wsHub  *service.Hub
}

var _ boards.UseCase = &implUsecase{}

func New(l log.Logger, repo repository.Repository, userUC user.UseCase, roleUC role.UseCase, wsHub *service.Hub) boards.UseCase {
	return &implUsecase{
		l:      l,
		clock:  util.Now,
		repo:   repo,
		userUC: userUC,
		roleUC: roleUC,
		wsHub:  wsHub,
	}
}
