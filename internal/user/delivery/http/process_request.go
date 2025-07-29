package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	pkgErrors "gitlab.com/tantai-kanban/kanban-api/pkg/errors"
	"gitlab.com/tantai-kanban/kanban-api/pkg/scope"
)

func (h handler) processDetailRequest(c *gin.Context) (string, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.user.delivery.http.processDetailRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return "", models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	id := c.Param("id")
	if id == "" {
		h.l.Errorf(ctx, "internal.user.delivery.http.processDetailRequest.c.Param: missing id parameter")
		return "", models.Scope{}, errWrongQuery
	}

	return id, scope.NewScope(p), nil
}

func (h handler) processDetailMeRequest(c *gin.Context) (models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.user.delivery.http.processDetailMeRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	return scope.NewScope(p), nil
}

func (h handler) processUpdateProfileRequest(c *gin.Context) (updateProfileReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.user.delivery.http.processUpdateProfileRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return updateProfileReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req updateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.user.delivery.http.processUpdateProfileRequest.c.ShouldBindJSON: %v", err)
		return updateProfileReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processCreateRequest(c *gin.Context) (createReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.user.delivery.http.processCreateRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return createReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.user.delivery.http.processCreateRequest.c.ShouldBindJSON: %v", err)
		return createReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}
