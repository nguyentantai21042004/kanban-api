package http

import (
	"github.com/nguyentantai21042004/kanban-api/internal/upload"
	"github.com/nguyentantai21042004/kanban-api/pkg/discord"
	pkgLog "github.com/nguyentantai21042004/kanban-api/pkg/log"
)

type handler struct {
	l  pkgLog.Logger
	uc upload.UseCase
	d  *discord.Discord
}

func New(l pkgLog.Logger, uc upload.UseCase, d *discord.Discord) Handler {
	return &handler{l: l, uc: uc, d: d}
}
