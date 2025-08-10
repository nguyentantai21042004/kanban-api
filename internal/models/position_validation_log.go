package models

import (
	"encoding/json"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
)

type PositionValidationLog struct {
	ID                 string                 `json:"id"`
	ListID             *string                `json:"list_id,omitempty"`
	BoardID            *string                `json:"board_id,omitempty"`
	ValidationType     string                 `json:"validation_type"`
	TargetType         string                 `json:"target_type"`
	IsValid            bool                   `json:"is_valid"`
	ErrorMessage       *string                `json:"error_message,omitempty"`
	ErrorDetails       map[string]interface{} `json:"error_details,omitempty"`
	RecordCount        *int                   `json:"record_count,omitempty"`
	ProblematicRecords map[string]interface{} `json:"problematic_records,omitempty"`
	CheckedAt          *time.Time             `json:"checked_at,omitempty"`
}

func NewPositionValidationLog(dbPositionValidationLog dbmodels.PositionValidationLog) PositionValidationLog {
	var errorDetails map[string]interface{}
	if dbPositionValidationLog.ErrorDetails.Valid {
		_ = json.Unmarshal(dbPositionValidationLog.ErrorDetails.JSON, &errorDetails)
	}

	var problematicRecords map[string]interface{}
	if dbPositionValidationLog.ProblematicRecords.Valid {
		_ = json.Unmarshal(dbPositionValidationLog.ProblematicRecords.JSON, &problematicRecords)
	}

	return PositionValidationLog{
		ID:                 dbPositionValidationLog.ID,
		ListID:             dbPositionValidationLog.ListID.Ptr(),
		BoardID:            dbPositionValidationLog.BoardID.Ptr(),
		ValidationType:     dbPositionValidationLog.ValidationType,
		TargetType:         dbPositionValidationLog.TargetType,
		IsValid:            dbPositionValidationLog.IsValid,
		ErrorMessage:       dbPositionValidationLog.ErrorMessage.Ptr(),
		ErrorDetails:       errorDetails,
		RecordCount:        dbPositionValidationLog.RecordCount.Ptr(),
		ProblematicRecords: problematicRecords,
		CheckedAt:          dbPositionValidationLog.CheckedAt.Ptr(),
	}
}
