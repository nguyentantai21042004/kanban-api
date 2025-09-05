package postgres

import (
	"database/sql"
	"time"

	"github.com/nguyentantai21042004/kanban-api/internal/lists/repository"
	"github.com/nguyentantai21042004/kanban-api/pkg/log"
	"github.com/nguyentantai21042004/kanban-api/pkg/util"
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
