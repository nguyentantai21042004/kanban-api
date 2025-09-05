package usecase

import (
	"time"

	"github.com/nguyentantai21042004/kanban-api/internal/cards"
	"github.com/nguyentantai21042004/kanban-api/internal/comments"
	"github.com/nguyentantai21042004/kanban-api/internal/comments/repository"
	"github.com/nguyentantai21042004/kanban-api/internal/user"
	"github.com/nguyentantai21042004/kanban-api/internal/websocket/service"
	"github.com/nguyentantai21042004/kanban-api/pkg/log"
	"github.com/nguyentantai21042004/kanban-api/pkg/util"
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
