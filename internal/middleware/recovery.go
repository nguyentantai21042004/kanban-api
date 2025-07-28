package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/pkg/discord"
	"gitlab.com/tantai-kanban/kanban-api/pkg/response"
)

func Recovery(d *discord.Discord) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[Panic Recovered] %v\n", err)
				log.Printf("Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
				response.PanicError(c, err, d)
			}
		}()
		c.Next()
	}
}
