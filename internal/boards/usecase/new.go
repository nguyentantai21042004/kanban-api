package usecase

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/boards"
	"gitlab.com/tantai-kanban/kanban-api/internal/boards/repository"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

type implUsecase struct {
	l     log.Logger
	repo  repository.Repository
	clock func() time.Time
}

var _ boards.UseCase = &implUsecase{}

func New(l log.Logger, repo repository.Repository) boards.UseCase {
	return &implUsecase{
		l:     l,
		repo:  repo,
		clock: util.Now,
	}
}
