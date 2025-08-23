package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/internal/admin"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	pkgErrors "gitlab.com/tantai-kanban/kanban-api/pkg/errors"
	"gitlab.com/tantai-kanban/kanban-api/pkg/scope"
)

type dashboardReq struct {
	Period string `form:"period"`
}

func (h handler) processDashboardRequest(c *gin.Context) (dashboardReq, models.Scope, error) {
	ctx := c.Request.Context()

	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "internal.admin.delivery.http.processDashboardRequest.jwt.GetPayloadFromContext: %v", "payload not found")
		return dashboardReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req dashboardReq
	if err := c.ShouldBindQuery(&req); err != nil {
		return dashboardReq{}, models.Scope{}, err
	}

	return req, scope.NewScope(p), nil
}

func (r dashboardReq) toInput() admin.DashboardInput {
	return admin.DashboardInput{Period: r.Period}
}

// Admin users list/create/update

type usersReq struct {
	Search string `form:"search"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

func (h handler) processUsersRequest(c *gin.Context) (usersReq, models.Scope, error) {
	ctx := c.Request.Context()
	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		return usersReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}
	var req usersReq
	if err := c.ShouldBindQuery(&req); err != nil {
		return usersReq{}, models.Scope{}, err
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	return req, scope.NewScope(p), nil
}

func (r usersReq) toInput() admin.UsersInput {
	return admin.UsersInput{Search: r.Search, Page: r.Page, PerPage: r.Limit}
}

type createUserReq struct {
	Username string `json:"username" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	RoleID   string `json:"role_id"`
	Password string `json:"password"`
}

func (r createUserReq) toInput() admin.CreateUserInput {
	return admin.CreateUserInput{Username: r.Username, FullName: r.FullName, RoleID: r.RoleID, Password: r.Password}
}

type updateUserReq struct {
	FullName  *string `json:"full_name"`
	RoleID    *string `json:"role_id"`
	RoleAlias *string `json:"role_alias"`
	IsActive  *bool   `json:"is_active"`
}

func (r updateUserReq) toInput() admin.UpdateUserInput {
	return admin.UpdateUserInput{FullName: r.FullName, RoleID: r.RoleID, RoleAlias: r.RoleAlias, IsActive: r.IsActive}
}

func (h handler) processCreateUserRequest(c *gin.Context) (createUserReq, models.Scope, error) {
	ctx := c.Request.Context()
	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		return createUserReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}
	var req createUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return createUserReq{}, models.Scope{}, err
	}
	return req, scope.NewScope(p), nil
}

func (h handler) processUpdateUserRequest(c *gin.Context) (string, models.Scope, updateUserReq, error) {
	ctx := c.Request.Context()
	p, ok := scope.GetPayloadFromContext(ctx)
	if !ok {
		return "", models.Scope{}, updateUserReq{}, pkgErrors.NewUnauthorizedHTTPError()
	}
	id := c.Param("id")
	var req updateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return "", models.Scope{}, updateUserReq{}, err
	}
	return id, scope.NewScope(p), req, nil
}
