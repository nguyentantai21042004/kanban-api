package http

import (
	"slices"

	"github.com/nguyentantai21042004/kanban-api/internal/role"
	"github.com/nguyentantai21042004/kanban-api/pkg/response"

	"github.com/gin-gonic/gin"
)

// @Summary Get role detail
// @Description Get a role by ID
// @Tags Role
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param id path string true "Role ID"
// @Success 200 {object} roleItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/roles/{id} [GET]
func (h handler) Detail(c *gin.Context) {
	ctx := c.Request.Context()

	id, sc, err := h.processDetailRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.role.http.Detail.processDetailRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Detail(ctx, sc, id)
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.role.http.Detail.uc.Detail: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.role.http.Detail.uc.Detail: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary List all roles
// @Description Get list of all roles in system
// @Tags Role
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Success 200 {object} []roleItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/roles [GET]
func (h handler) List(c *gin.Context) {
	ctx := c.Request.Context()

	sc, err := h.processListRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.role.http.List.processListRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	roles, err := h.uc.List(ctx, sc, role.ListInput{})
	if err != nil {
		mapErr := h.mapErrorCode(err)
		h.l.Errorf(ctx, "internal.role.http.List.uc.List: %v", err)
		response.Error(c, mapErr, h.d)
		return
	}

	items := make([]roleItem, len(roles))
	for i, r := range roles {
		items[i] = h.newItem(r)
	}

	response.OK(c, items)
}
