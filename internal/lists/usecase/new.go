package usecase

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/boards"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/websocket/service"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/position"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

type implUsecase struct {
	l          log.Logger
	repo       repository.Repository
	wsHub      *service.Hub
	positionUC position.Usecase
	boardUC    boards.UseCase
	clock      func() time.Time
}

var _ lists.UseCase = &implUsecase{}

func New(l log.Logger, repo repository.Repository, wsHub *service.Hub, positionUC position.Usecase, boardUC boards.UseCase) lists.UseCase {
	return &implUsecase{
		l:          l,
		repo:       repo,
		wsHub:      wsHub,
		positionUC: positionUC,
		boardUC:    boardUC,
		clock:      util.Now,
	}
}
