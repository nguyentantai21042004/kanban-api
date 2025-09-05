package models

import (
	"time"

	"github.com/nguyentantai21042004/kanban-api/internal/dbmodels"
)

type PositionStatistic struct {
	ID                string     `json:"id"`
	ListID            *string    `json:"list_id,omitempty"`
	BoardID           *string    `json:"board_id,omitempty"`
	TargetType        string     `json:"target_type"`
	RecordCount       int        `json:"record_count"`
	AvgLength         *float64   `json:"avg_length,omitempty"`
	MaxLength         *int       `json:"max_length,omitempty"`
	MinLength         *int       `json:"min_length,omitempty"`
	LengthStddev      *float64   `json:"length_stddev,omitempty"`
	LongKeyCount      *int       `json:"long_key_count,omitempty"`
	LongKeyPercentage *float64   `json:"long_key_percentage,omitempty"`
	NeedsRebalance    *bool      `json:"needs_rebalance,omitempty"`
	HealthScore       *float64   `json:"health_score,omitempty"`
	PerformanceImpact *string    `json:"performance_impact,omitempty"`
	CalculatedAt      *time.Time `json:"calculated_at,omitempty"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
}

func NewPositionStatistic(dbPositionStatistic dbmodels.PositionStatistic) PositionStatistic {
	var avgLength *float64
	if dbPositionStatistic.AvgLength.Big != nil {
		f, _ := dbPositionStatistic.AvgLength.Big.Float64()
		avgLength = &f
	}

	var lengthStddev *float64
	if dbPositionStatistic.LengthStddev.Big != nil {
		f, _ := dbPositionStatistic.LengthStddev.Big.Float64()
		lengthStddev = &f
	}

	var longKeyPercentage *float64
	if dbPositionStatistic.LongKeyPercentage.Big != nil {
		f, _ := dbPositionStatistic.LongKeyPercentage.Big.Float64()
		longKeyPercentage = &f
	}

	var healthScore *float64
	if dbPositionStatistic.HealthScore.Big != nil {
		f, _ := dbPositionStatistic.HealthScore.Big.Float64()
		healthScore = &f
	}

	return PositionStatistic{
		ID:                dbPositionStatistic.ID,
		ListID:            dbPositionStatistic.ListID.Ptr(),
		BoardID:           dbPositionStatistic.BoardID.Ptr(),
		TargetType:        dbPositionStatistic.TargetType,
		RecordCount:       dbPositionStatistic.RecordCount,
		AvgLength:         avgLength,
		MaxLength:         dbPositionStatistic.MaxLength.Ptr(),
		MinLength:         dbPositionStatistic.MinLength.Ptr(),
		LengthStddev:      lengthStddev,
		LongKeyCount:      dbPositionStatistic.LongKeyCount.Ptr(),
		LongKeyPercentage: longKeyPercentage,
		NeedsRebalance:    dbPositionStatistic.NeedsRebalance.Ptr(),
		HealthScore:       healthScore,
		PerformanceImpact: dbPositionStatistic.PerformanceImpact.Ptr(),
		CalculatedAt:      dbPositionStatistic.CalculatedAt.Ptr(),
		ExpiresAt:         dbPositionStatistic.ExpiresAt.Ptr(),
	}
}
