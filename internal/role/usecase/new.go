package usecase

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/role"
	"gitlab.com/tantai-kanban/kanban-api/internal/role/repository"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

type usecase struct {
	l    pkgLog.Logger
	repo repository.Repository
}

func New(l pkgLog.Logger, repo repository.Repository) role.UseCase {
	return &usecase{
		l:    l,
		repo: repo,
	}
}
