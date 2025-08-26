package usecase

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/boards"
	"gitlab.com/tantai-kanban/kanban-api/internal/cards"
	"gitlab.com/tantai-kanban/kanban-api/internal/cards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists"
	"gitlab.com/tantai-kanban/kanban-api/internal/role"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
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
	listUC     lists.UseCase
	userUC     user.UseCase
	roleUC     role.UseCase
	clock      func() time.Time
}

var _ cards.UseCase = &implUsecase{}

func New(l log.Logger, repo repository.Repository, wsHub *service.Hub, positionUC position.Usecase, boardUC boards.UseCase, listUC lists.UseCase, userUC user.UseCase, roleUC role.UseCase) cards.UseCase {
	return &implUsecase{
		l:          l,
		repo:       repo,
		wsHub:      wsHub,
		positionUC: positionUC,
		clock:      util.Now,
		boardUC:    boardUC,
		listUC:     listUC,
		userUC:     userUC,
		roleUC:     roleUC,
	}
}
