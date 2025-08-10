package models

import (
	"encoding/json"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
)

type MigrationProgress struct {
	ID                string                 `json:"id"`
	TableName         string                 `json:"table_name"`
	TargetType        string                 `json:"target_type"`
	TotalRecords      int64                  `json:"total_records"`
	MigratedRecords   *int64                 `json:"migrated_records,omitempty"`
	FailedRecords     *int64                 `json:"failed_records,omitempty"`
	StartedAt         *time.Time             `json:"started_at,omitempty"`
	CompletedAt       *time.Time             `json:"completed_at,omitempty"`
	LastUpdatedAt     *time.Time             `json:"last_updated_at,omitempty"`
	Status            *string                `json:"status,omitempty"`
	ErrorDetails      map[string]interface{} `json:"error_details,omitempty"`
	MigrationStrategy *string                `json:"migration_strategy,omitempty"`
	CreatedAt         *time.Time             `json:"created_at,omitempty"`
}

func NewMigrationProgress(dbMigrationProgress dbmodels.MigrationProgress) MigrationProgress {
	var errorDetails map[string]interface{}
	if dbMigrationProgress.ErrorDetails.Valid {
		_ = json.Unmarshal(dbMigrationProgress.ErrorDetails.JSON, &errorDetails)
	}

	return MigrationProgress{
		ID:                dbMigrationProgress.ID,
		TableName:         dbMigrationProgress.TableName,
		TargetType:        dbMigrationProgress.TargetType,
		TotalRecords:      dbMigrationProgress.TotalRecords,
		MigratedRecords:   dbMigrationProgress.MigratedRecords.Ptr(),
		FailedRecords:     dbMigrationProgress.FailedRecords.Ptr(),
		StartedAt:         dbMigrationProgress.StartedAt.Ptr(),
		CompletedAt:       dbMigrationProgress.CompletedAt.Ptr(),
		LastUpdatedAt:     dbMigrationProgress.LastUpdatedAt.Ptr(),
		Status:            dbMigrationProgress.Status.Ptr(),
		ErrorDetails:      errorDetails,
		MigrationStrategy: dbMigrationProgress.MigrationStrategy.Ptr(),
		CreatedAt:         dbMigrationProgress.CreatedAt.Ptr(),
	}
}
