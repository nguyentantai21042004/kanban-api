package http

import (
	"slices"

	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/pkg/response"
)

// @Summary Get cards
// @Description Get cards
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDc2MzUwODcsImp0aSI6IjIwMjUtMDUtMTIgMTM6MTE6MjcuODI5ODQ0NTUxICswNzAwICswNyBtPSszNS4zNTAzNTUxMTAiLCJuYmYiOjE3NDcwMzAyODcsInN1YiI6ImM0NTk2MzAzLWRlNDItNDI0Yi1hZmNiLWVhNWJlNjNhYjA2MCIsImVtYWlsIjoidGFpMjEwNDIwMDRAZ21haWwuY29tIiwidHlwZSI6ImFjY2VzcyIsInJlZnJlc2giOmZhbHNlfQ.NxH8MvILhwWo02PDybh8ofJpz8rnSA71EO6lwZs3ykQ)
// @Param ids query string false "IDs"
// @Param list_id query string false "List ID"
// @Param keyword query string false "Keyword"
// @Param page query integer false "Page"
// @Param limit query integer false "Limit"
// @Success 200 {object} getCardResp "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards [GET]
func (h handler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processGetRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.Get.processGetRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Get(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.Get.uc.Get: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.Get.uc.Get: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newGetResp(o))
}

// @Summary Create card
// @Description Create a new card
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body createReq true "Card data"
// @Success 200 {object} cardItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards [POST]
func (h handler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processCreateRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.Create.processCreateRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Create(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.Create.uc.Create: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.Create.uc.Create: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Update card
// @Description Update an existing card
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body updateReq true "Card data"
// @Success 200 {object} cardItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards [PUT]
func (h handler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processUpdateRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.Update.processUpdateRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Update(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.Update.uc.Update: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.Update.uc.Update: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Get card detail
// @Description Get a card by ID
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param id path string true "Card ID"
// @Success 200 {object} cardItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards/{id} [GET]
func (h handler) Detail(c *gin.Context) {
	ctx := c.Request.Context()

	id, sc, err := h.processDetailRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.Detail.processDetailRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Detail(ctx, sc, id)
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.Detail.uc.Detail: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.Detail.uc.Detail: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Delete card
// @Description Delete a card by ID
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body deleteReq true "Card IDs"
// @Success 200 {object} response.Resp "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards [DELETE]
func (h handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processDeleteRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.Delete.processDeleteRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	err = h.uc.Delete(ctx, sc, req.IDs)
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.Delete.uc.Delete: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.Delete.uc.Delete: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, nil)
}

// @Summary Move card
// @Description Move a card to different list/position
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body moveReq true "Move data"
// @Success 200 {object} cardItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards/move [POST]
func (h handler) Move(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processMoveRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.Move.processMoveRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Move(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.Move.uc.Move: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.Move.uc.Move: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Get card activities
// @Description Get activities for a card
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param card_id query string true "Card ID"
// @Param page query integer false "Page"
// @Param limit query integer false "Limit"
// @Success 200 {object} getCardActivitiesResp "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards/activities [GET]
func (h handler) GetActivities(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processGetActivitiesRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.GetActivities.processGetActivitiesRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.GetActivities(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.GetActivities.uc.GetActivities: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.GetActivities.uc.GetActivities: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newGetActivitiesResp(o))
}
