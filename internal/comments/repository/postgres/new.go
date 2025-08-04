package postgres

import (
	"database/sql"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/comments/repository"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

type implRepository struct {
	l        log.Logger
	database *sql.DB
	clock    func() time.Time
}

var _ repository.Repository = implRepository{}

func New(l log.Logger, database *sql.DB) implRepository {
	return implRepository{
		l:        l,
		database: database,
		clock:    util.Now,
	}
}
