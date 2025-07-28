package usecase

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/labels"
	"gitlab.com/tantai-kanban/kanban-api/internal/labels/repository"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

type implUsecase struct {
	l     log.Logger
	repo  repository.Repository
	clock func() time.Time
}

var _ labels.UseCase = &implUsecase{}

func New(l log.Logger, repo repository.Repository) labels.UseCase {
	return &implUsecase{
		l:     l,
		repo:  repo,
		clock: util.Now,
	}
}
