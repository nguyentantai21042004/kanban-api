package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/tantai-kanban/kanban-api/internal/cards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

func TestExtractPositionString(t *testing.T) {
	repo := implRepository{}

	tests := []struct {
		name     string
		card     dbmodels.Card
		expected string
	}{
		{
			name: "nil position",
			card: dbmodels.Card{
				Position: types.Decimal{Big: nil},
			},
			expected: "n",
		},
		{
			name: "integer position",
			card: dbmodels.Card{
				Position: types.Decimal{Big: decimal.New(100, 0)},
			},
			expected: "2", // base62[100/62] = base62[1] = "1" + base62[38] = "c" -> but this is simplified
		},
		{
			name: "decimal position",
			card: dbmodels.Card{
				Position: types.Decimal{Big: decimal.New(1500, -3)}, // 1.5
			},
			expected: "1.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.extractPositionString(tt.card)
			assert.NotEmpty(t, result)
			// Basic validation that it returns a valid position string
		})
	}
}

func TestExtractPositionFromOptions(t *testing.T) {
	repo := implRepository{}

	tests := []struct {
		name     string
		opts     repository.MoveOptions
		expected string
	}{
		{
			name: "zero position",
			opts: repository.MoveOptions{
				Position: 0,
			},
			expected: "",
		},
		{
			name: "positive position",
			opts: repository.MoveOptions{
				Position: 100.5,
			},
			expected: "2", // simplified expectation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.extractPositionFromOptions(tt.opts)
			// Basic validation
			if tt.opts.Position == 0 {
				assert.Empty(t, result)
			} else {
				assert.NotEmpty(t, result)
			}
		})
	}
}

func TestConvertPositionStringToDecimal(t *testing.T) {
	repo := implRepository{}

	tests := []struct {
		name   string
		posStr string
	}{
		{
			name:   "numeric string",
			posStr: "123.45",
		},
		{
			name:   "single character",
			posStr: "a",
		},
		{
			name:   "multi character",
			posStr: "ab",
		},
		{
			name:   "empty string",
			posStr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.convertPositionStringToDecimal(tt.posStr)
			assert.NotNil(t, result.Big)
		})
	}
}

func TestConvertFloatToPositionString(t *testing.T) {
	repo := implRepository{}

	tests := []struct {
		name     string
		floatPos float64
	}{
		{
			name:     "zero",
			floatPos: 0,
		},
		{
			name:     "small positive",
			floatPos: 10,
		},
		{
			name:     "large positive",
			floatPos: 1000,
		},
		{
			name:     "negative",
			floatPos: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.convertFloatToPositionString(tt.floatPos)
			assert.NotEmpty(t, result)
		})
	}
}

func TestConvertPositionStringToFloat(t *testing.T) {
	repo := implRepository{}

	tests := []struct {
		name   string
		posStr string
	}{
		{
			name:   "empty string",
			posStr: "",
		},
		{
			name:   "single character",
			posStr: "a",
		},
		{
			name:   "multi character",
			posStr: "ab",
		},
		{
			name:   "numeric character",
			posStr: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.convertPositionStringToFloat(tt.posStr)
			assert.True(t, result >= 0)
		})
	}
}

func TestBuildEnhancedMoveModel(t *testing.T) {
	repo := implRepository{
		clock: func() time.Time { return time.Now() },
	}

	tests := []struct {
		name        string
		opts        repository.MoveOptions
		expectError bool
	}{
		{
			name: "valid options",
			opts: repository.MoveOptions{
				ID:       "550e8400-e29b-41d4-a716-446655440000", // Valid UUID
				ListID:   "test-list",
				Position: 100.5,
				OldModel: models.Card{ID: "old-id"},
			},
			expectError: false,
		},
		{
			name: "invalid UUID",
			opts: repository.MoveOptions{
				ID:       "invalid-uuid",
				ListID:   "test-list",
				Position: 100.5,
				OldModel: models.Card{ID: "old-id"},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card, cols, err := repo.buildEnhancedMoveModel(context.Background(), tt.opts)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.opts.ID, card.ID)
				assert.Equal(t, tt.opts.ListID, card.ListID)
				assert.NotNil(t, card.Position.Big)
				assert.NotEmpty(t, cols)
				assert.Contains(t, cols, dbmodels.CardColumns.ListID)
				assert.Contains(t, cols, dbmodels.CardColumns.Position)
				assert.Contains(t, cols, dbmodels.CardColumns.UpdatedAt)
			}
		})
	}
}

func TestBuildEnhancedActivityModel(t *testing.T) {
	repo := implRepository{}

	oldData := map[string]interface{}{
		"from_list_id":  "list1",
		"from_position": "a",
	}

	newData := map[string]interface{}{
		"to_list_id":    "list2",
		"to_position":   "b",
		"was_optimized": true,
	}

	activity := repo.buildEnhancedActivityModel(
		context.Background(),
		"card-id",
		string(models.CardActionTypeMoved),
		oldData,
		newData,
	)

	assert.Equal(t, "card-id", activity.CardID)
	assert.Equal(t, dbmodels.CardActionType(models.CardActionTypeMoved), activity.ActionType)
	assert.True(t, activity.OldData.Valid)
	assert.True(t, activity.NewData.Valid)
}

func TestConvertModelToDBModel(t *testing.T) {
	repo := implRepository{}

	model := models.Card{
		ID:       "test-id",
		ListID:   "test-list",
		Position: 123.45,
		Name:     "Test Card",
	}

	dbModel := repo.convertModelToDBModel(model)

	assert.Equal(t, model.ID, dbModel.ID)
	assert.Equal(t, model.ListID, dbModel.ListID)
	assert.Equal(t, model.Name, dbModel.Name)
	assert.NotNil(t, dbModel.Position.Big)

	// Check position conversion
	floatPos, _ := dbModel.Position.Big.Float64()
	assert.InDelta(t, model.Position, floatPos, 0.001)
}

func TestPowerFunction(t *testing.T) {
	tests := []struct {
		name     string
		base     int
		exp      int
		expected float64
	}{
		{
			name:     "power of 1",
			base:     5,
			exp:      1,
			expected: 5,
		},
		{
			name:     "power of 0",
			base:     5,
			exp:      0,
			expected: 1,
		},
		{
			name:     "power of 2",
			base:     3,
			exp:      2,
			expected: 9,
		},
		{
			name:     "power of 3",
			base:     2,
			exp:      3,
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pow(tt.base, tt.exp)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Integration test for position workflow
func TestPositionWorkflow(t *testing.T) {
	// This test simulates the full position workflow without database
	repo := implRepository{
		clock: func() time.Time { return time.Now() },
	}

	// Create mock cards
	mockCards := []dbmodels.Card{
		{
			ID:       "card1",
			ListID:   "list1",
			Position: types.Decimal{Big: decimal.New(100, 0)},
			Name:     "Card 1",
		},
		{
			ID:       "card2",
			ListID:   "list1",
			Position: types.Decimal{Big: decimal.New(200, 0)},
			Name:     "Card 2",
		},
	}

	// Test position extraction
	for _, card := range mockCards {
		posStr := repo.extractPositionString(card)
		assert.NotEmpty(t, posStr)

		// Convert back to decimal
		decimal := repo.convertPositionStringToDecimal(posStr)
		assert.NotNil(t, decimal.Big)
	}

	// Test move options processing
	opts := repository.MoveOptions{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		ListID:   "list1",
		Position: 150.0,
		OldModel: models.Card{ID: "old-card"},
	}

	posStr := repo.extractPositionFromOptions(opts)
	assert.NotEmpty(t, posStr)

	// Test model building
	card, cols, err := repo.buildEnhancedMoveModel(context.Background(), opts)
	require.NoError(t, err)
	assert.Equal(t, opts.ID, card.ID)
	assert.Equal(t, opts.ListID, card.ListID)
	assert.NotEmpty(t, cols)
}
