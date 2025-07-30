package usecase

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
	"gitlab.com/tantai-kanban/kanban-api/internal/user/repository"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
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
