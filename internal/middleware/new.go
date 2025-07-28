package middleware

import (
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
	pkgScope "gitlab.com/tantai-kanban/kanban-api/pkg/scope"
)

type Middleware struct {
	l          pkgLog.Logger
	jwtManager pkgScope.Manager
}

func New(l pkgLog.Logger, jwtManager pkgScope.Manager) Middleware {
	return Middleware{
		l:          l,
		jwtManager: jwtManager,
	}
}
