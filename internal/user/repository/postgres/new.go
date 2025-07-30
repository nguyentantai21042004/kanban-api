package postgres

import (
	"database/sql"

	"gitlab.com/tantai-kanban/kanban-api/internal/user/repository"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

type implRepository struct {
	l        pkgLog.Logger
	database *sql.DB
}

var _ repository.Repository = &implRepository{}

func New(l pkgLog.Logger, database *sql.DB) *implRepository {
	return &implRepository{
		l:        l,
		database: database,
	}
}
