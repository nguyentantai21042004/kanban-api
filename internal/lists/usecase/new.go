package usecase

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/lists"
	"gitlab.com/tantai-kanban/kanban-api/internal/lists/repository"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

type implUsecase struct {
	l     log.Logger
	repo  repository.Repository
	clock func() time.Time
}

var _ lists.UseCase = &implUsecase{}

func New(l log.Logger, repo repository.Repository) lists.UseCase {
	return &implUsecase{
		l:     l,
		repo:  repo,
		clock: util.Now,
	}
}
