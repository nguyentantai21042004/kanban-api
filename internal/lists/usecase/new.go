package usecase

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/lists"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/websocket/service"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

type implUsecase struct {
	l     log.Logger
	repo  repository.Repository
	wsHub *service.Hub
	clock func() time.Time
}

var _ lists.UseCase = &implUsecase{}

func New(l log.Logger, repo repository.Repository, wsHub *service.Hub) lists.UseCase {
	return &implUsecase{
		l:     l,
		repo:  repo,
		wsHub: wsHub,
		clock: util.Now,
	}
}
