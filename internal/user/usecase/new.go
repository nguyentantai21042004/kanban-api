package usecase

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

type usecase struct {
	l    log.Logger
	repo user.Repository
}

func New(l log.Logger, repo user.Repository) user.UseCase {
	return &usecase{
		l:    l,
		repo: repo,
	}
}
