package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyentantai21042004/kanban-api/pkg/response"
)

// @Summary Get user detail
// @Description Get user detail by ID (Super Admin only or own profile)
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param id path string true "User ID"
// @Success 200 {object} userItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/users/{id} [GET]
func (h handler) Detail(c *gin.Context) {
	ctx := c.Request.Context()

	id, sc, err := h.processDetailRequest(c)
	if err != nil {
		response.Error(c, h.mapErrorCode(err), h.d)
		return
	}

	o, err := h.uc.Detail(ctx, sc, id)
	if err != nil {
		h.l.Errorf(ctx, "internal.user.http.Detail.uc.Detail: %v", err)
		response.Error(c, h.mapErrorCode(err), h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Get my profile
// @Description Get current user profile
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {object} userItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/users/me [GET]
func (h handler) DetailMe(c *gin.Context) {
	ctx := c.Request.Context()

	sc, err := h.processDetailMeRequest(c)
	if err != nil {
		response.Error(c, h.mapErrorCode(err), h.d)
		return
	}

	o, err := h.uc.DetailMe(ctx, sc)
	if err != nil {
		h.l.Errorf(ctx, "internal.user.http.DetailMe.uc.DetailMe: %v", err)
		response.Error(c, h.mapErrorCode(err), h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Update profile
// @Description Update current user profile
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param request body updateProfileReq true "Update profile request"
// @Success 200 {object} userItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/users/profile [PUT]
func (h handler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processUpdateProfileRequest(c)
	if err != nil {
		response.Error(c, h.mapErrorCode(err), h.d)
		return
	}

	o, err := h.uc.UpdateProfile(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "internal.user.http.UpdateProfile.uc.UpdateProfile: %v", err)
		response.Error(c, h.mapErrorCode(err), h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Create user
// @Description Create new user (Super Admin only)
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param request body createReq true "Create user request"
// @Success 201 {object} userItem "Created"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/users [POST]
func (h handler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processCreateRequest(c)
	if err != nil {
		response.Error(c, h.mapErrorCode(err), h.d)
		return
	}

	o, err := h.uc.Create(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "internal.user.http.Create.uc.Create: %v", err)
		response.Error(c, h.mapErrorCode(err), h.d)
		return
	}

	response.OK(c, h.newItem(o))
}
