package models

import (
	"encoding/json"
	"time"

	"github.com/nguyentantai21042004/kanban-api/internal/dbmodels"
)

type RebalanceJob struct {
	ID            string                 `json:"id"`
	ListID        *string                `json:"list_id,omitempty"`
	BoardID       *string                `json:"board_id,omitempty"`
	TargetType    string                 `json:"target_type"`
	Priority      *string                `json:"priority,omitempty"`
	Status        *string                `json:"status,omitempty"`
	TriggerReason *string                `json:"trigger_reason,omitempty"`
	Strategy      *string                `json:"strategy,omitempty"`
	ScheduledAt   *time.Time             `json:"scheduled_at,omitempty"`
	StartedAt     *time.Time             `json:"started_at,omitempty"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty"`
	ErrorMessage  *string                `json:"error_message,omitempty"`
	Result        map[string]interface{} `json:"result,omitempty"`
	Attempts      *int                   `json:"attempts,omitempty"`
	MaxAttempts   *int                   `json:"max_attempts,omitempty"`
	CreatedAt     *time.Time             `json:"created_at,omitempty"`
	CreatedBy     *string                `json:"created_by,omitempty"`
}

func NewRebalanceJob(dbRebalanceJob dbmodels.RebalanceJob) RebalanceJob {
	var result map[string]interface{}
	if dbRebalanceJob.Result.Valid {
		_ = json.Unmarshal(dbRebalanceJob.Result.JSON, &result)
	}

	return RebalanceJob{
		ID:            dbRebalanceJob.ID,
		ListID:        dbRebalanceJob.ListID.Ptr(),
		BoardID:       dbRebalanceJob.BoardID.Ptr(),
		TargetType:    dbRebalanceJob.TargetType,
		Priority:      dbRebalanceJob.Priority.Ptr(),
		Status:        dbRebalanceJob.Status.Ptr(),
		TriggerReason: dbRebalanceJob.TriggerReason.Ptr(),
		Strategy:      dbRebalanceJob.Strategy.Ptr(),
		ScheduledAt:   dbRebalanceJob.ScheduledAt.Ptr(),
		StartedAt:     dbRebalanceJob.StartedAt.Ptr(),
		CompletedAt:   dbRebalanceJob.CompletedAt.Ptr(),
		ErrorMessage:  dbRebalanceJob.ErrorMessage.Ptr(),
		Result:        result,
		Attempts:      dbRebalanceJob.Attempts.Ptr(),
		MaxAttempts:   dbRebalanceJob.MaxAttempts.Ptr(),
		CreatedAt:     dbRebalanceJob.CreatedAt.Ptr(),
		CreatedBy:     dbRebalanceJob.CreatedBy.Ptr(),
	}
}
