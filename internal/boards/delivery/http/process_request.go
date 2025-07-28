package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	pkgErrors "gitlab.com/tantai-kanban/kanban-api/pkg/errors"
	"gitlab.com/tantai-kanban/kanban-api/pkg/scope"
)

func (h handler) processGetRequest(c *gin.Context) (getReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.boards.delivery.http.processGetRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return getReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req getReq
	if err := c.ShouldBindQuery(&req); err != nil {
		h.l.Errorf(ctx, "internal.boards.delivery.http.processGetRequest.c.ShouldBindQuery: %v", err)
		return getReq{}, models.Scope{}, errWrongQuery
	}

	req.PageQuery.Adjust()
	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "internal.boards.delivery.http.processGetRequest.req.validate: %v", err)
		return getReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}
