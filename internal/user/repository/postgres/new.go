package postgre

import (
	"database/sql"

	"gitlab.com/tantai-kanban/kanban-api/internal/user"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

type repository struct {
	l        pkgLog.Logger
	database *sql.DB
}

func New(l pkgLog.Logger, database *sql.DB) user.Repository {
	return &repository{
		l:        l,
		database: database,
	}
}
