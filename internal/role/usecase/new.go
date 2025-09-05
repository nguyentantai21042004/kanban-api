package usecase

import (
	"github.com/nguyentantai21042004/kanban-api/internal/role"
	"github.com/nguyentantai21042004/kanban-api/internal/role/repository"
	pkgLog "github.com/nguyentantai21042004/kanban-api/pkg/log"
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
