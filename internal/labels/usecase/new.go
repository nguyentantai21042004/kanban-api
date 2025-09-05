package usecase

import (
	"time"

	"github.com/nguyentantai21042004/kanban-api/internal/labels"
	"github.com/nguyentantai21042004/kanban-api/internal/labels/repository"
	"github.com/nguyentantai21042004/kanban-api/pkg/log"
	"github.com/nguyentantai21042004/kanban-api/pkg/util"
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
