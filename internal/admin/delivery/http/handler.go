package http

import (
	"slices"

	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/pkg/response"
)

// @Summary Admin dashboard
// @Description KPI overview and activity chart
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param period query string false "7d | 30d | 90d" default(7d)
// @Success 200 {object} dashboardResp "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/admin/dashboard [GET]
func (h handler) Dashboard(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processDashboardRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.admin.http.Dashboard.processDashboardRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Dashboard(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.admin.http.Dashboard.uc.Dashboard: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.admin.http.Dashboard.uc.Dashboard: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	respondOK(c, newDashboardResp(o))
}

// @Summary Admin system health
// @Description Overall system status for monitoring tab
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Success 200 {object} admin.HealthOutput "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/admin/health [GET]
func (h handler) Health(c *gin.Context) {
	ctx := c.Request.Context()
	// reusing scope extraction not needed here since protected by Auth middleware, but uc.Health may need scope
	_, sc, err := h.processDashboardRequest(c)
	if err != nil {
		response.Error(c, err, h.d)
		return
	}
	o, err := h.uc.Health(ctx, sc)
	if err != nil {
		response.Error(c, err, h.d)
		return
	}
	respondOK(c, o)
}

// @Summary Admin list users
// @Description List users for admin management
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param search query string false "Search by username/full_name"
// @Param page query integer false "Page" default(1)
// @Param limit query integer false "Limit" default(20)
// @Success 200 {object} usersResp "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/admin/users [GET]
func (h handler) Users(c *gin.Context) {
	ctx := c.Request.Context()
	req, sc, err := h.processUsersRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.admin.http.Users.processUsersRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}
	o, err := h.uc.Users(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "internal.admin.http.Users.uc.Users: %v", err)
		response.Error(c, err, h.d)
		return
	}
	respondOK(c, newUsersResp(o))
}

// @Summary Admin create user
// @Description Create new user for admin
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body createUserReq true "Create user"
// @Success 200 {object} admin.UserItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/admin/users [POST]
func (h handler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	req, sc, err := h.processCreateUserRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.admin.http.CreateUser.processCreateUserRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}
	o, err := h.uc.CreateUser(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "internal.admin.http.CreateUser.uc.CreateUser: %v", err)
		response.Error(c, err, h.d)
		return
	}
	respondOK(c, o)
}

// @Summary Admin update user
// @Description Quick update user info for admin
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param id path string true "User ID"
// @Param body body updateUserReq true "Update user"
// @Success 200 {object} admin.UserItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/admin/users/{id} [PUT]
func (h handler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	id, sc, req, err := h.processUpdateUserRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.admin.http.UpdateUser.processUpdateUserRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}
	o, err := h.uc.UpdateUser(ctx, sc, id, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "internal.admin.http.UpdateUser.uc.UpdateUser: %v", err)
		response.Error(c, err, h.d)
		return
	}
	respondOK(c, o)
}

// @Summary Admin list roles
// @Description List all roles for admin to select when creating users
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Success 200 {object} []admin.RoleItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/admin/roles [GET]
func (h handler) Roles(c *gin.Context) {
	ctx := c.Request.Context()
	_, sc, err := h.processDashboardRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.admin.http.Roles.processDashboardRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}
	roles, err := h.uc.Roles(ctx, sc)
	if err != nil {
		h.l.Errorf(ctx, "internal.admin.http.Roles.uc.Roles: %v", err)
		response.Error(c, err, h.d)
		return
	}
	respondOK(c, roles)
}
