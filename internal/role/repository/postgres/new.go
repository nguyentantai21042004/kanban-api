package postgres

import (
	"database/sql"

	"gitlab.com/tantai-kanban/kanban-api/internal/role/repository"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

type implRepository struct {
	l        pkgLog.Logger
	database *sql.DB
}

var _ repository.Repository = &implRepository{}

func New(l pkgLog.Logger, database *sql.DB) repository.Repository {
	return &implRepository{
		l:        l,
		database: database,
	}
}
