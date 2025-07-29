package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	pkgErrors "gitlab.com/tantai-kanban/kanban-api/pkg/errors"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
	"gitlab.com/tantai-kanban/kanban-api/pkg/scope"
)

func (h handler) processGetRequest(c *gin.Context) (getReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processGetRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return getReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req getReq
	if err := c.ShouldBindQuery(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processGetRequest.c.ShouldBindQuery: %v", err)
		return getReq{}, models.Scope{}, errWrongQuery
	}

	req.PageQuery.Adjust()
	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processGetRequest.req.validate: %v", err)
		return getReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processCreateRequest(c *gin.Context) (createReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processCreateRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return createReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processCreateRequest.c.ShouldBindQuery: %v", err)
		return createReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processUpdateRequest(c *gin.Context) (updateReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processUpdateRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return updateReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req updateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processUpdateRequest.c.ShouldBindQuery: %v", err)
		return updateReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processDetailRequest(c *gin.Context) (string, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processDetailRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return "", models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	id := c.Param("id")
	if err := postgres.IsUUID(id); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processDetailRequest.c.Param: %v", err)
		return "", models.Scope{}, errWrongQuery
	}

	return id, scope.NewScope(p), nil
}

func (h handler) processDeleteRequest(c *gin.Context) (deleteReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processDetailRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return deleteReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req deleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processDeleteRequest.c.ShouldBindQuery: %v", err)
		return deleteReq{}, models.Scope{}, errWrongQuery
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processDeleteRequest.req.validate: %v", err)
		return deleteReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processMoveRequest(c *gin.Context) (moveReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processMoveRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return moveReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req moveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processMoveRequest.c.ShouldBindJSON: %v", err)
		return moveReq{}, models.Scope{}, errWrongQuery
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processMoveRequest.req.validate: %v", err)
		return moveReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processGetActivitiesRequest(c *gin.Context) (getActivitiesReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processGetActivitiesRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return getActivitiesReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req getActivitiesReq
	if err := c.ShouldBindQuery(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processGetActivitiesRequest.c.ShouldBindQuery: %v", err)
		return getActivitiesReq{}, models.Scope{}, errWrongQuery
	}

	req.PageQuery.Adjust()
	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processGetActivitiesRequest.req.validate: %v", err)
		return getActivitiesReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}
