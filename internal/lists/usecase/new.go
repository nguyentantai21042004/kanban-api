package usecase

import (
	"time"

	"github.com/nguyentantai21042004/kanban-api/internal/boards"
	"github.com/nguyentantai21042004/kanban-api/internal/lists"
	"github.com/nguyentantai21042004/kanban-api/internal/lists/repository"
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
	userUC     user.UseCase
	roleUC     role.UseCase
	clock      func() time.Time
}

var _ lists.UseCase = &implUsecase{}

func New(l log.Logger, repo repository.Repository, wsHub *service.Hub, positionUC position.Usecase, boardUC boards.UseCase, userUC user.UseCase, roleUC role.UseCase) lists.UseCase {
	return &implUsecase{
		l:          l,
		repo:       repo,
		wsHub:      wsHub,
		positionUC: positionUC,
		boardUC:    boardUC,
		userUC:     userUC,
		roleUC:     roleUC,
		clock:      util.Now,
	}
}
