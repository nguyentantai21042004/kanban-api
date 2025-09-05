package usecase

import (
	"github.com/nguyentantai21042004/kanban-api/internal/user"
	"github.com/nguyentantai21042004/kanban-api/internal/user/repository"
	"github.com/nguyentantai21042004/kanban-api/pkg/log"
)

type usecase struct {
	l    log.Logger
	repo repository.Repository
}

func New(l log.Logger, repo repository.Repository) user.UseCase {
	return &usecase{
		l:    l,
		repo: repo,
	}
}
