package http

import "github.com/gin-gonic/gin"

type Handler interface {
	Detail(c *gin.Context)
	DetailMe(c *gin.Context)
	UpdateProfile(c *gin.Context)
	Create(c *gin.Context) // Chỉ Super Admin
}
