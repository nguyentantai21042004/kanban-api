package postgres

import (
	"database/sql"

	"gitlab.com/tantai-kanban/kanban-api/internal/upload"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

type repository struct {
	l        pkgLog.Logger
	database *sql.DB
}

func New(l pkgLog.Logger, database *sql.DB) upload.Repository {
	return &repository{
		l:        l,
		database: database,
	}
}
