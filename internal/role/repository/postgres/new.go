package postgres

import (
	"database/sql"

	"github.com/nguyentantai21042004/kanban-api/internal/role/repository"
	pkgLog "github.com/nguyentantai21042004/kanban-api/pkg/log"
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
