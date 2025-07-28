package smtp

import (
	"gitlab.com/tantai-kanban/kanban-api/config"
	"gitlab.com/tantai-kanban/kanban-api/internal/core/smtp"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

type implService struct {
	l   log.Logger
	cfg config.SMTPConfig
}

func New(l log.Logger, cfg config.SMTPConfig) smtp.UseCase {
	return implService{
		l:   l,
		cfg: cfg,
	}
}
