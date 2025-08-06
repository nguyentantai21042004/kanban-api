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
// @Param board_id query string false "Board ID"
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

// Enhanced functionality methods

// @Summary Assign card
// @Description Assign a card to a user
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body assignReq true "Assignment data"
// @Success 200 {object} cardItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards/assign [POST]
func (h handler) Assign(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processAssignRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.Assign.processAssignRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Assign(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.Assign.uc.Assign: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.Assign.uc.Assign: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Unassign card
// @Description Remove assignment from a card
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body unassignReq true "Unassignment data"
// @Success 200 {object} cardItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards/unassign [POST]
func (h handler) Unassign(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processUnassignRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.Unassign.processUnassignRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.Unassign(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.Unassign.uc.Unassign: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.Unassign.uc.Unassign: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Add attachment to card
// @Description Add an attachment to a card
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body addAttachmentReq true "Attachment data"
// @Success 200 {object} cardItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards/attachments/add [POST]
func (h handler) AddAttachment(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processAddAttachmentRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.AddAttachment.processAddAttachmentRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.AddAttachment(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.AddAttachment.uc.AddAttachment: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.AddAttachment.uc.AddAttachment: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Remove attachment from card
// @Description Remove an attachment from a card
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body removeAttachmentReq true "Attachment data"
// @Success 200 {object} cardItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards/attachments/remove [POST]
func (h handler) RemoveAttachment(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processRemoveAttachmentRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.RemoveAttachment.processRemoveAttachmentRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.RemoveAttachment(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.RemoveAttachment.uc.RemoveAttachment: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.RemoveAttachment.uc.RemoveAttachment: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Update time tracking
// @Description Update time tracking for a card
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body updateTimeTrackingReq true "Time tracking data"
// @Success 200 {object} cardItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards/time-tracking [PUT]
func (h handler) UpdateTimeTracking(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processUpdateTimeTrackingRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.UpdateTimeTracking.processUpdateTimeTrackingRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.UpdateTimeTracking(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.UpdateTimeTracking.uc.UpdateTimeTracking: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.UpdateTimeTracking.uc.UpdateTimeTracking: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Update checklist
// @Description Update checklist for a card
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body updateChecklistReq true "Checklist data"
// @Success 200 {object} cardItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards/checklist [PUT]
func (h handler) UpdateChecklist(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processUpdateChecklistRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.UpdateChecklist.processUpdateChecklistRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.UpdateChecklist(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.UpdateChecklist.uc.UpdateChecklist: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.UpdateChecklist.uc.UpdateChecklist: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Add tag to card
// @Description Add a tag to a card
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body addTagReq true "Tag data"
// @Success 200 {object} cardItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards/tags/add [POST]
func (h handler) AddTag(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processAddTagRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.AddTag.processAddTagRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.AddTag(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.AddTag.uc.AddTag: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.AddTag.uc.AddTag: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Remove tag from card
// @Description Remove a tag from a card
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body removeTagReq true "Tag data"
// @Success 200 {object} cardItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards/tags/remove [POST]
func (h handler) RemoveTag(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processRemoveTagRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.RemoveTag.processRemoveTagRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.RemoveTag(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.RemoveTag.uc.RemoveTag: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.RemoveTag.uc.RemoveTag: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Set start date
// @Description Set start date for a card
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body setStartDateReq true "Start date data"
// @Success 200 {object} cardItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards/start-date [PUT]
func (h handler) SetStartDate(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processSetStartDateRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.SetStartDate.processSetStartDateRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.SetStartDate(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.SetStartDate.uc.SetStartDate: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.SetStartDate.uc.SetStartDate: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}

// @Summary Set completion date
// @Description Set completion date for a card
// @Tags Card
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer <token>)
// @Param body body setCompletionDateReq true "Completion date data"
// @Success 200 {object} cardItem "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 404 {object} response.Resp "Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/cards/completion-date [PUT]
func (h handler) SetCompletionDate(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processSetCompletionDateRequest(c)
	if err != nil {
		h.l.Warnf(ctx, "internal.cards.http.SetCompletionDate.processSetCompletionDateRequest: %v", err)
		response.Error(c, err, h.d)
		return
	}

	o, err := h.uc.SetCompletionDate(ctx, sc, req.toInput())
	if err != nil {
		mapErr := h.mapErrorCode(err)
		if slices.Contains(NotFound, mapErr) {
			h.l.Warnf(ctx, "internal.cards.http.SetCompletionDate.uc.SetCompletionDate: %v", err)
		} else {
			h.l.Errorf(ctx, "internal.cards.http.SetCompletionDate.uc.SetCompletionDate: %v", err)
		}
		response.Error(c, mapErr, h.d)
		return
	}

	response.OK(c, h.newItem(o))
}
