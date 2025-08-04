package http

import (
	"slices"

	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/pkg/response"
)

// @Summary Get comments
// @Description Get comments
// @Tags Comment
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDc2MzUwODcsImp0aSI6IjIwMjUtMDUtMTIgMTM6MTE6MjcuODI5ODQ0NTUxICswNzAwICswNyBtPSszNS4zNTAzNTUxMTAiLCJuYmYiOjE3NDcwMzAyODcsInN1YiI6ImM0NTk2MzAzLWRlNDItNDI0Yi1hZmNiLWVhNWJlNjNhYjA2MCIsImVtYWlsIjoidGFpMjEwNDIwMDRAZ21haWwuY29tIiwidHlwZSI6ImFjY2VzcyIsInJlZnJlc2giOmZhbHNlfQ.NxH8MvILhwWo02PDybh8ofJpz8rnSA71EO6lwZs3ykQ)
// @Param ids query string false "IDs"
// @Param keyword query string false "Keyword"
// @Param card_id query string false "Card ID"
// @Param user_id query string false "User ID"
// @Param parent_id query string false "Parent ID"
// @Param page query integer false "Page"
// @Param limit query integer false "Limit"
// @Success 200 {object} getCommentResp "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/comments [GET]
func (h handler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processGetRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.comments.http.Get.processGetRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Get(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.comments.http.Get.uc.Get: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.comments.http.Get.uc.Get: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newGetResp(o))
}

// @Summary Create comment
// @Description Create a new comment
// @Tags Comment
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body createReq true "Comment data"
// @Success 200 {object} commentItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/comments [POST]
func (h handler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processCreateRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.comments.http.Create.processCreateRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Create(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.comments.http.Create.uc.Create: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.comments.http.Create.uc.Create: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Update comment
// @Description Update an existing comment
// @Tags Comment
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param id path string true "Comment ID"
// @Param body body updateReq true "Comment data"
// @Success 200 {object} commentItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/comments/{id} [PUT]
func (h handler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processUpdateRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.comments.http.Update.processUpdateRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Update(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.comments.http.Update.uc.Update: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.comments.http.Update.uc.Update: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Get comment detail
// @Description Get comment detail by ID
// @Tags Comment
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param id path string true "Comment ID"
// @Success 200 {object} commentItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/comments/{id} [GET]
func (h handler) Detail(c *gin.Context) {
	ctx := c.Request.Context()

	id, sc, err := h.processDetailRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.comments.http.Detail.processDetailRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Detail(ctx, sc, id)
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.comments.http.Detail.uc.Detail: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.comments.http.Detail.uc.Detail: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Delete comments
// @Description Delete comments by IDs
// @Tags Comment
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param ids query string true "Comment IDs (comma separated)"
// @Success 200 {object} response.Resp "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/comments [DELETE]
func (h handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processDeleteRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.comments.http.Delete.processDeleteRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	err = h.uc.Delete(ctx, sc, req.IDs)
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.comments.http.Delete.uc.Delete: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.comments.http.Delete.uc.Delete: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, nil)
}

// @Summary Get comments by card
// @Description Get comments by card ID
// @Tags Comment
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param card_id path string true "Card ID"
// @Success 200 {object} getCommentResp "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards/{card_id}/comments [GET]
func (h handler) GetByCard(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processGetByCardRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.comments.http.GetByCard.processGetByCardRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.GetByCard(ctx, sc, req.CardID)
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.comments.http.GetByCard.uc.GetByCard: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.comments.http.GetByCard.uc.GetByCard: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newGetResp(o))
}
