package http

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/role"
	"gitlab.com/tantai-kanban/kanban-api/pkg/discord"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

type handler struct {
	l  pkgLog.Logger
	uc role.UseCase
	d  *discord.Discord
}

func New(
	l pkgLog.Logger,
	uc role.UseCase,
	d *discord.Discord,
) Handler {
	return &handler{
		l:  l,
		uc: uc,
		d:  d,
	}
}
