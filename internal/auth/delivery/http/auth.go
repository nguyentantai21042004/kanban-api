package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyentantai21042004/kanban-api/pkg/response"
)

// @Summary Login
// @Description Login with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body loginReq true "Login request"
// @Success 200 {object} loginResp "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/auth/login [POST]
func (h handler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processLoginRequest(c)
	if err != nil {
		response.Error(c, h.mapErrorCode(err), h.d)
		return
	}

	o, err := h.uc.Login(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "internal.auth.http.Login.uc.Login: %v", err)
		response.Error(c, h.mapErrorCode(err), h.d)
		return
	}

	response.OK(c, h.newLoginResp(o))
}

// @Summary Refresh token
// @Description Refresh access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body refreshTokenReq true "Refresh token request"
// @Success 200 {object} refreshTokenResp "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/auth/refresh [POST]
func (h handler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processRefreshTokenRequest(c)
	if err != nil {
		response.Error(c, h.mapErrorCode(err), h.d)
		return
	}

	o, err := h.uc.RefreshToken(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "internal.auth.http.RefreshToken.uc.RefreshToken: %v", err)
		response.Error(c, h.mapErrorCode(err), h.d)
		return
	}

	response.OK(c, h.newRefreshTokenResp(o))
}

// @Summary Logout
// @Description Logout and invalidate tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {object} response.Resp "Success"
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/auth/logout [POST]
func (h handler) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	sc, err := h.processLogoutRequest(c)
	if err != nil {
		response.Error(c, h.mapErrorCode(err), h.d)
		return
	}

	err = h.uc.Logout(ctx, sc)
	if err != nil {
		h.l.Errorf(ctx, "internal.auth.http.Logout.uc.Logout: %v", err)
		response.Error(c, h.mapErrorCode(err), h.d)
		return
	}

	response.OK(c, nil)
}
