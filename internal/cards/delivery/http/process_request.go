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

// Enhanced functionality process methods

func (h handler) processAssignRequest(c *gin.Context) (assignReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processAssignRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return assignReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req assignReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processAssignRequest.c.ShouldBindJSON: %v", err)
		return assignReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processUnassignRequest(c *gin.Context) (unassignReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processUnassignRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return unassignReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req unassignReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processUnassignRequest.c.ShouldBindJSON: %v", err)
		return unassignReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processAddAttachmentRequest(c *gin.Context) (addAttachmentReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processAddAttachmentRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return addAttachmentReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req addAttachmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processAddAttachmentRequest.c.ShouldBindJSON: %v", err)
		return addAttachmentReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processRemoveAttachmentRequest(c *gin.Context) (removeAttachmentReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processRemoveAttachmentRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return removeAttachmentReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req removeAttachmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processRemoveAttachmentRequest.c.ShouldBindJSON: %v", err)
		return removeAttachmentReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processUpdateTimeTrackingRequest(c *gin.Context) (updateTimeTrackingReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processUpdateTimeTrackingRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return updateTimeTrackingReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req updateTimeTrackingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processUpdateTimeTrackingRequest.c.ShouldBindJSON: %v", err)
		return updateTimeTrackingReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processUpdateChecklistRequest(c *gin.Context) (updateChecklistReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processUpdateChecklistRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return updateChecklistReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req updateChecklistReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processUpdateChecklistRequest.c.ShouldBindJSON: %v", err)
		return updateChecklistReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processAddTagRequest(c *gin.Context) (addTagReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processAddTagRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return addTagReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req addTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processAddTagRequest.c.ShouldBindJSON: %v", err)
		return addTagReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processRemoveTagRequest(c *gin.Context) (removeTagReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processRemoveTagRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return removeTagReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req removeTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processRemoveTagRequest.c.ShouldBindJSON: %v", err)
		return removeTagReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processSetStartDateRequest(c *gin.Context) (setStartDateReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processSetStartDateRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return setStartDateReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req setStartDateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processSetStartDateRequest.c.ShouldBindJSON: %v", err)
		return setStartDateReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}

func (h handler) processSetCompletionDateRequest(c *gin.Context) (setCompletionDateReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processSetCompletionDateRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return setCompletionDateReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req setCompletionDateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.cards.delivery.http.processSetCompletionDateRequest.c.ShouldBindJSON: %v", err)
		return setCompletionDateReq{}, models.Scope{}, errWrongQuery
	}

	return req, scope.NewScope(p), nil
}
