package usecase

import (
	"context"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/admin"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/metrics"
)

func (uc implUsecase) Health(ctx context.Context, sc models.Scope) (admin.HealthOutput, error) {
	// Basic API status and timestamps
	now := time.Now().UTC()

	// WS connections: count all clients across boards
	wsConns := 0
	if uc.wsHub != nil {
		boards, _ := uc.wsHub.GetConnectedBoards(ctx)
		for _, b := range boards {
			if n, err := uc.wsHub.GetActiveUsersCount(ctx, b); err == nil {
				wsConns += n
			}
		}
	}

	avgMs, uptime := metrics.SnapshotHTTP()
	out := admin.HealthOutput{
		APIStatus:               "healthy",
		ResponseTimeMs:          avgMs,
		UptimePercentage:        uptime,
		WebsocketConnections:    wsConns,
		WebsocketMessagesPerSec: 0,
		WebsocketAvgLatencyMs:   0,
		CheckedAt:               now.Format(time.RFC3339),
	}
	return out, nil
}
