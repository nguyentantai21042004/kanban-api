package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	pkgErrors "github.com/nguyentantai21042004/kanban-api/pkg/errors"
	"github.com/nguyentantai21042004/kanban-api/pkg/scope"
)

func (h handler) processLoginRequest(c *gin.Context) (loginReq, models.Scope, error) {
	ctx := c.Request.Context()

	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.auth.delivery.http.processLoginRequest.c.ShouldBindJSON: %v", err)
		return loginReq{}, models.Scope{}, errWrongQuery
	}

	return req, models.Scope{}, nil
}

func (h handler) processRefreshTokenRequest(c *gin.Context) (refreshTokenReq, models.Scope, error) {
	ctx := c.Request.Context()

	var req refreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "internal.auth.delivery.http.processRefreshTokenRequest.c.ShouldBindJSON: %v", err)
		return refreshTokenReq{}, models.Scope{}, errWrongQuery
	}

	return req, models.Scope{}, nil
}

func (h handler) processLogoutRequest(c *gin.Context) (models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.auth.delivery.http.processLogoutRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	return scope.NewScope(p), nil
}
