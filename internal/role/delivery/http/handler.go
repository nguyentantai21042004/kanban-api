package http

import (
	"slices"

	"gitlab.com/tantai-kanban/kanban-api/pkg/response"

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
