package http

import (
	"slices"

	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/pkg/response"
)

// @Summary Get labels
// @Description Get labels
// @Tags Label
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDc2MzUwODcsImp0aSI6IjIwMjUtMDUtMTIgMTM6MTE6MjcuODI5ODQ0NTUxICswNzAwICswNyBtPSszNS4zNTAzNTUxMTAiLCJuYmYiOjE3NDcwMzAyODcsInN1YiI6ImM0NTk2MzAzLWRlNDItNDI0Yi1hZmNiLWVhNWJlNjNhYjA2MCIsImVtYWlsIjoidGFpMjEwNDIwMDRAZ21haWwuY29tIiwidHlwZSI6ImFjY2VzcyIsInJlZnJlc2giOmZhbHNlfQ.NxH8MvILhwWo02PDybh8ofJpz8rnSA71EO6lwZs3ykQ)
// @Param ids query string false "IDs"
// @Param label_id query string false "label ID"
// @Param keyword query string false "Keyword"
// @Param page query integer false "Page"
// @Param limit query integer false "Limit"
// @Success 200 {object} getLabelResp "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/labels [GET]
func (h handler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processGetRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.labels.http.Get.processGetRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Get(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.labels.http.Get.uc.Get: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.labels.http.Get.uc.Get: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newGetResp(o))
}

// @Summary Create label
// @Description Create a new label
// @Tags Label
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body createReq true "Label data"
// @Success 200 {object} labelItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/labels [POST]
func (h handler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processCreateRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.labels.http.Create.processCreateRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Create(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.labels.http.Create.uc.Create: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.labels.http.Create.uc.Create: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Update label
// @Description Update an existing label
// @Tags label
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body updateReq true "Label data"
// @Success 200 {object} labelItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/labels [PUT]
func (h handler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processUpdateRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.labels.http.Update.processUpdateRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Update(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.labels.http.Update.uc.Update: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.labels.http.Update.uc.Update: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Get label detail
// @Description Get a label by ID
// @Tags Label
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param id path string true "Label ID"
// @Success 200 {object} labelItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/labels/{id} [GET]
func (h handler) Detail(c *gin.Context) {
	ctx := c.Request.Context()

	id, sc, err := h.processDetailRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.labels.http.Detail.processDetailRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Detail(ctx, sc, id)
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.labels.http.Detail.uc.Detail: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.labels.http.Detail.uc.Detail: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Delete label
// @Description Delete a label by ID
// @Tags Label
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body deleteReq true "Label IDs"
// @Success 200 {object} response.Resp "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/labels [DELETE]
func (h handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processDeleteRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.labels.http.Delete.processDeleteRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	err = h.uc.Delete(ctx, sc, req.IDs)
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.labels.http.Delete.uc.Delete: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.labels.http.Delete.uc.Delete: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, nil)
}
