package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	pkgErrors "github.com/nguyentantai21042004/kanban-api/pkg/errors"
	"github.com/nguyentantai21042004/kanban-api/pkg/postgres"
	"github.com/nguyentantai21042004/kanban-api/pkg/scope"
)

func (h handler) processGetRequest(c *gin.Context) (getReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.labels.delivery.http.processGetRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return getReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req getReq
	if err := c.ShouldBindQuery(&req); err != nil {
		h.l.Errorf(ctx, "internal.labels.delivery.http.processGetRequest.c.ShouldBindQuery: %v", err)
		return getReq{}, models.Scope{}, errWrongQuery
	}

	req.PageQuery.Adjust()
	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "internal.labels.delivery.http.processGetRequest.req.validate: %v", err)
		return getReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processCreateRequest(c *gin.Context) (createReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.labels.delivery.http.processCreateRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return createReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.labels.delivery.http.processCreateRequest.c.ShouldBindQuery: %v", err)
		return createReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processUpdateRequest(c *gin.Context) (updateReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.labels.delivery.http.processUpdateRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return updateReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req updateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.labels.delivery.http.processUpdateRequest.c.ShouldBindQuery: %v", err)
		return updateReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processDetailRequest(c *gin.Context) (string, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.labels.delivery.http.processDetailRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return "", models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	id := c.Param("id")
	if err := postgres.IsUUID(id); err != nil {
		h.l.Errorf(ctx, "internal.labels.delivery.http.processDetailRequest.c.Param: %v", err)
		return "", models.Scope{}, errWrongQuery
	}

	return id, scope.NewScope(p), nil
}

func (h handler) processDeleteRequest(c *gin.Context) (deleteReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.labels.delivery.http.processDetailRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return deleteReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req deleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.labels.delivery.http.processDeleteRequest.c.ShouldBindQuery: %v", err)
		return deleteReq{}, models.Scope{}, errWrongQuery
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "internal.boards.delivery.http.processDeleteRequest.req.validate: %v", err)
		return deleteReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}
