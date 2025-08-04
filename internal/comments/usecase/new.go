package usecase

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/cards"
	"gitlab.com/tantai-kanban/kanban-api/internal/comments"
	"gitlab.com/tantai-kanban/kanban-api/internal/comments/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
	"gitlab.com/tantai-kanban/kanban-api/internal/websocket/service"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

type implUsecase struct {
	l       log.Logger
	clock   func() time.Time
	repo    repository.Repository
	userUC  user.UseCase
	cardsUC cards.UseCase
	wsHub   *service.Hub
}

var _ comments.UseCase = &implUsecase{}

func New(l log.Logger, repo repository.Repository, userUC user.UseCase, cardsUC cards.UseCase, wsHub *service.Hub) comments.UseCase {
	return &implUsecase{
		l:       l,
		clock:   util.Now,
		repo:    repo,
		userUC:  userUC,
		cardsUC: cardsUC,
		wsHub:   wsHub,
	}
}
