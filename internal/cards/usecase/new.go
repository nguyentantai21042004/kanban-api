package usecase

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/cards"
	"gitlab.com/tantai-kanban/kanban-api/internal/cards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
	"gitlab.com/tantai-kanban/kanban-api/internal/websocket/service"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

type implUsecase struct {
	l      log.Logger
	repo   repository.Repository
	wsHub  *service.Hub
	clock  func() time.Time
	listUC lists.UseCase
	userUC user.UseCase
}

var _ cards.UseCase = &implUsecase{}

func New(l log.Logger, repo repository.Repository, wsHub *service.Hub, listUC lists.UseCase, userUC user.UseCase) cards.UseCase {
	return &implUsecase{
		l:      l,
		repo:   repo,
		wsHub:  wsHub,
		clock:  util.Now,
		listUC: listUC,
		userUC: userUC,
	}
}
