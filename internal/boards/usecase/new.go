package usecase

import (
	"time"

	"github.com/nguyentantai21042004/kanban-api/internal/boards"
	"github.com/nguyentantai21042004/kanban-api/internal/boards/repository"
	"github.com/nguyentantai21042004/kanban-api/internal/lists"
	"github.com/nguyentantai21042004/kanban-api/internal/role"
	"github.com/nguyentantai21042004/kanban-api/internal/user"
	"github.com/nguyentantai21042004/kanban-api/internal/websocket/service"
	"github.com/nguyentantai21042004/kanban-api/pkg/log"
	"github.com/nguyentantai21042004/kanban-api/pkg/util"
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
