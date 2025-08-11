package models

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
)

type RebalanceEvent struct {
	ID              string     `json:"id"`
	ListID          *string    `json:"list_id,omitempty"`
	BoardID         *string    `json:"board_id,omitempty"`
	TargetType      string     `json:"target_type"`
	Strategy        *string    `json:"strategy,omitempty"`
	RecordCount     *int       `json:"record_count,omitempty"`
	AvgLengthBefore *float64   `json:"avg_length_before,omitempty"`
	AvgLengthAfter  *float64   `json:"avg_length_after,omitempty"`
	MaxLengthBefore *int       `json:"max_length_before,omitempty"`
	MaxLengthAfter  *int       `json:"max_length_after,omitempty"`
	MinLengthBefore *int       `json:"min_length_before,omitempty"`
	MinLengthAfter  *int       `json:"min_length_after,omitempty"`
	DurationMS      *int       `json:"duration_ms,omitempty"`
	TriggerReason   *string    `json:"trigger_reason,omitempty"`
	JobID           *string    `json:"job_id,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
}

func NewRebalanceEvent(dbRebalanceEvent dbmodels.RebalanceEvent) RebalanceEvent {
	var avgLengthBefore *float64
	if dbRebalanceEvent.AvgLengthBefore.Big != nil {
		f, _ := dbRebalanceEvent.AvgLengthBefore.Big.Float64()
		avgLengthBefore = &f
	}

	var avgLengthAfter *float64
	if dbRebalanceEvent.AvgLengthAfter.Big != nil {
		f, _ := dbRebalanceEvent.AvgLengthAfter.Big.Float64()
		avgLengthAfter = &f
	}

	return RebalanceEvent{
		ID:              dbRebalanceEvent.ID,
		ListID:          dbRebalanceEvent.ListID.Ptr(),
		BoardID:         dbRebalanceEvent.BoardID.Ptr(),
		TargetType:      dbRebalanceEvent.TargetType,
		Strategy:        dbRebalanceEvent.Strategy.Ptr(),
		RecordCount:     dbRebalanceEvent.RecordCount.Ptr(),
		AvgLengthBefore: avgLengthBefore,
		AvgLengthAfter:  avgLengthAfter,
		MaxLengthBefore: dbRebalanceEvent.MaxLengthBefore.Ptr(),
		MaxLengthAfter:  dbRebalanceEvent.MaxLengthAfter.Ptr(),
		MinLengthBefore: dbRebalanceEvent.MinLengthBefore.Ptr(),
		MinLengthAfter:  dbRebalanceEvent.MinLengthAfter.Ptr(),
		DurationMS:      dbRebalanceEvent.DurationMS.Ptr(),
		TriggerReason:   dbRebalanceEvent.TriggerReason.Ptr(),
		JobID:           dbRebalanceEvent.JobID.Ptr(),
		CreatedAt:       dbRebalanceEvent.CreatedAt.Ptr(),
	}
}
