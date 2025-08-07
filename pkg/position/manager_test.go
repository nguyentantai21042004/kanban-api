package position

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()
	assert.NotNil(t, manager)
	assert.Equal(t, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", manager.base62Chars)
	assert.Equal(t, "n", manager.midPoint)
}

func TestGeneratePosition(t *testing.T) {
	manager := NewManager()

	tests := []struct {
		name     string
		before   string
		after    string
		expected string
		hasError bool
	}{
		{
			name:     "empty before and after",
			before:   "",
			after:    "",
			expected: "n",
			hasError: false,
		},
		{
			name:     "empty before",
			before:   "",
			after:    "z",
			expected: "y",
			hasError: false,
		},
		{
			name:     "empty after",
			before:   "a",
			after:    "",
			expected: "b",
			hasError: false,
		},
		{
			name:     "normal between",
			before:   "a",
			after:    "c",
			expected: "b",
			hasError: false,
		},
		{
			name:     "invalid order",
			before:   "z",
			after:    "a",
			expected: "",
			hasError: true,
		},
		{
			name:     "adjacent characters",
			before:   "a",
			after:    "b",
			expected: "an",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := manager.GeneratePosition(tt.before, tt.after)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestValidateAndFixPosition(t *testing.T) {
	manager := NewManager()

	// Test cards in order
	cards := []Card{
		{ID: "1", ListID: "list1", Position: "a", Name: "Card 1"},
		{ID: "2", ListID: "list1", Position: "b", Name: "Card 2"},
		{ID: "3", ListID: "list1", Position: "c", Name: "Card 3"},
	}

	tests := []struct {
		name              string
		cardID            string
		targetListID      string
		requestedPosition string
		expectedWasFixed  bool
		expectedError     bool
	}{
		{
			name:              "valid position",
			cardID:            "4",
			targetListID:      "list1",
			requestedPosition: "d",
			expectedWasFixed:  false,
			expectedError:     false,
		},
		{
			name:              "empty position",
			cardID:            "4",
			targetListID:      "list1",
			requestedPosition: "",
			expectedWasFixed:  true,
			expectedError:     false,
		},
		{
			name:              "invalid position",
			cardID:            "4",
			targetListID:      "list1",
			requestedPosition: "!",
			expectedWasFixed:  true,
			expectedError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, wasFixed, err := manager.ValidateAndFixPosition(
				tt.cardID,
				tt.targetListID,
				tt.requestedPosition,
				cards,
			)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
				assert.Equal(t, tt.expectedWasFixed, wasFixed)
			}
		})
	}
}

func TestRebalancePositions(t *testing.T) {
	manager := NewManager()

	// Create cards with very long positions that need rebalancing
	cards := []Card{
		{ID: "1", ListID: "list1", Position: "aaaaaaaaaa", Name: "Card 1"}, // 10 chars
		{ID: "2", ListID: "list1", Position: "bbbbbbbbbb", Name: "Card 2"}, // 10 chars
		{ID: "3", ListID: "list1", Position: "cccccccccc", Name: "Card 3"}, // 10 chars
	}

	rebalanceMap, err := manager.RebalancePositions(cards, 8)
	require.NoError(t, err)
	assert.Equal(t, 3, len(rebalanceMap))

	// Check that all cards got new positions
	for _, card := range cards {
		newPos, exists := rebalanceMap[card.ID]
		assert.True(t, exists)
		assert.NotEmpty(t, newPos)
		assert.True(t, len(newPos) <= 8, "New position should be shorter than max length")
	}
}

func TestBatchGeneratePositions(t *testing.T) {
	manager := NewManager()

	tests := []struct {
		name          string
		count         int
		before        string
		after         string
		expectedError bool
	}{
		{
			name:          "single position",
			count:         1,
			before:        "a",
			after:         "c",
			expectedError: false,
		},
		{
			name:          "multiple positions",
			count:         3,
			before:        "a",
			after:         "z",
			expectedError: false,
		},
		{
			name:          "zero count",
			count:         0,
			before:        "a",
			after:         "c",
			expectedError: true,
		},
		{
			name:          "negative count",
			count:         -1,
			before:        "a",
			after:         "c",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positions, err := manager.BatchGeneratePositions(tt.count, tt.before, tt.after)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, positions)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.count, len(positions))

				// Check that positions are in order
				for i := 1; i < len(positions); i++ {
					assert.True(t, positions[i-1] < positions[i], "Positions should be in ascending order")
				}
			}
		})
	}
}

func TestGetPositionMetrics(t *testing.T) {
	manager := NewManager()

	tests := []struct {
		name               string
		cards              []Card
		expectedTotalCards int
		expectedNeedsRebal bool
	}{
		{
			name:               "empty cards",
			cards:              []Card{},
			expectedTotalCards: 0,
			expectedNeedsRebal: false,
		},
		{
			name: "normal positions",
			cards: []Card{
				{ID: "1", Position: "a"},
				{ID: "2", Position: "b"},
				{ID: "3", Position: "c"},
			},
			expectedTotalCards: 3,
			expectedNeedsRebal: false,
		},
		{
			name: "long positions need rebalancing",
			cards: []Card{
				{ID: "1", Position: "aaaaaaaaaa"}, // 10 chars
				{ID: "2", Position: "bbbbbbbbbb"}, // 10 chars
			},
			expectedTotalCards: 2,
			expectedNeedsRebal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := manager.GetPositionMetrics(tt.cards)

			assert.Equal(t, tt.expectedTotalCards, metrics["total_cards"])
			assert.Equal(t, tt.expectedNeedsRebal, metrics["needs_rebalance"])

			if tt.expectedTotalCards > 0 {
				assert.NotNil(t, metrics["avg_length"])
				assert.NotNil(t, metrics["max_length"])
				assert.NotNil(t, metrics["min_length"])
			}
		})
	}
}

func TestComparePositions(t *testing.T) {
	manager := NewManager()

	tests := []struct {
		name     string
		a        string
		b        string
		expected int
	}{
		{
			name:     "a < b",
			a:        "a",
			b:        "b",
			expected: -1,
		},
		{
			name:     "a > b",
			a:        "b",
			b:        "a",
			expected: 1,
		},
		{
			name:     "a == b",
			a:        "a",
			b:        "a",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.ComparePositions(tt.a, tt.b)
			if tt.expected < 0 {
				assert.True(t, result < 0)
			} else if tt.expected > 0 {
				assert.True(t, result > 0)
			} else {
				assert.Equal(t, 0, result)
			}
		})
	}
}

// Integration test for position generation sequence
func TestPositionGenerationSequence(t *testing.T) {
	manager := NewManager()

	// Start with empty list
	positions := []string{}

	// Add 10 cards sequentially
	for i := 0; i < 10; i++ {
		var before, after string
		if len(positions) > 0 {
			before = positions[len(positions)-1] // Last position
		}

		newPos, err := manager.GeneratePosition(before, after)
		require.NoError(t, err)
		positions = append(positions, newPos)
	}

	// Verify all positions are in order
	for i := 1; i < len(positions); i++ {
		assert.True(t, positions[i-1] < positions[i],
			"Position %d (%s) should be less than position %d (%s)",
			i-1, positions[i-1], i, positions[i])
	}

	// Insert a card in the middle
	middlePos, err := manager.GeneratePosition(positions[4], positions[5])
	require.NoError(t, err)

	assert.True(t, positions[4] < middlePos)
	assert.True(t, middlePos < positions[5])
}
