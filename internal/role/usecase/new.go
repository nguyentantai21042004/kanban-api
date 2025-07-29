package usecase

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/role"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

type usecase struct {
	l    pkgLog.Logger
	repo role.Repository
}

func New(l pkgLog.Logger, repo role.Repository) role.UseCase {
	return &usecase{
		l:    l,
		repo: repo,
	}
}
