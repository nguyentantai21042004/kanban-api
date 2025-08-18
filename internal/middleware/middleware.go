package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/pkg/metrics"
	"gitlab.com/tantai-kanban/kanban-api/pkg/response"
	"gitlab.com/tantai-kanban/kanban-api/pkg/scope"
)

func (m Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		if tokenString == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		payload, err := m.jwtManager.Verify(tokenString)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		ctx = scope.SetPayloadToContext(ctx, payload)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// Metrics records simple per-request timing and uptime proxy
func (m Middleware) Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		status := c.Writer.Status()
		metrics.ObserveHTTP(status, duration)
	}
}
