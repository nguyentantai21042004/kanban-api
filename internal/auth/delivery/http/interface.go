package http

import "github.com/gin-gonic/gin"

type Handler interface {
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
	Logout(c *gin.Context)
}
