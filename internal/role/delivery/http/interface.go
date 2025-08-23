package http

import "github.com/gin-gonic/gin"

type Handler interface {
	Detail(c *gin.Context)
	List(c *gin.Context)
}
