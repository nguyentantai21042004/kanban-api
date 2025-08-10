package models

import (
	"encoding/json"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
)

type SystemConfig struct {
	Key         string                 `json:"key"`
	Value       map[string]interface{} `json:"value"`
	Description *string                `json:"description,omitempty"`
	UpdatedAt   *time.Time             `json:"updated_at,omitempty"`
	UpdatedBy   *string                `json:"updated_by,omitempty"`
}

func NewSystemConfig(dbSystemConfig dbmodels.SystemConfig) SystemConfig {
	var value map[string]interface{}
	if dbSystemConfig.Value.String() != "" {
		_ = json.Unmarshal([]byte(dbSystemConfig.Value.String()), &value)
	}

	return SystemConfig{
		Key:       dbSystemConfig.Key,
		Value:     value,
		UpdatedAt: dbSystemConfig.UpdatedAt.Ptr(),
	}
}
