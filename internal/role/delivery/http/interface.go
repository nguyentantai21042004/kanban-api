package http

import "github.com/gin-gonic/gin"

type Handler interface {
	Get(c *gin.Context)
	Detail(c *gin.Context)
}
