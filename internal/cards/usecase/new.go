package usecase

import (
	"time"

	"github.com/nguyentantai21042004/kanban-api/internal/boards"
	"github.com/nguyentantai21042004/kanban-api/internal/cards"
	"github.com/nguyentantai21042004/kanban-api/internal/cards/repository"
	"github.com/nguyentantai21042004/kanban-api/internal/lists"
	"github.com/nguyentantai21042004/kanban-api/internal/role"
	"github.com/nguyentantai21042004/kanban-api/internal/user"
	"github.com/nguyentantai21042004/kanban-api/internal/websocket/service"
	"github.com/nguyentantai21042004/kanban-api/pkg/log"
	"github.com/nguyentantai21042004/kanban-api/pkg/position"
	"github.com/nguyentantai21042004/kanban-api/pkg/util"
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
