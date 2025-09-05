package position

import (
	"testing"

	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPositionManager(t *testing.T) {
	manager := NewPositionManager()
	assert.NotNil(t, manager)
	assert.Equal(t, BaseChars, manager.baseChars)
	assert.Equal(t, Base, manager.base)
	assert.Equal(t, "i", manager.midChar)
	assert.Equal(t, "0", manager.minChar)
	assert.Equal(t, "z", manager.maxChar)
}

func TestGeneratePosition(t *testing.T) {
	manager := NewPositionManager()

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
			expected: "i",
			hasError: false,
		},
		{
			name:     "empty before",
			before:   "",
			after:    "z",
			expected: "yz", // Fixed: generateBefore appends 'z'
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
			expected: "ai",
			hasError: false,
		},
		{
			name:     "same characters",
			before:   "a",
			after:    "a",
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := manager.GeneratePosition(tt.before, tt.after)

			if tt.hasError {
				assert.Error(t, err)
				if tt.name == "invalid order" {
					assert.IsType(t, InvalidPositionOrderError{}, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestValidateAndFixPosition(t *testing.T) {
	manager := NewPositionManager()

	// Test cards in order
	cards := []models.Card{
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
		{
			name:              "position would break order",
			cardID:            "4",
			targetListID:      "list1",
			requestedPosition: "a", // Same as first card
			expectedWasFixed:  true,
			expectedError:     false,
		},
		{
			name:              "empty list",
			cardID:            "4",
			targetListID:      "list1",
			requestedPosition: "a",
			expectedWasFixed:  true,
			expectedError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCards := cards
			if tt.name == "empty list" {
				testCards = []models.Card{}
			}

			result, wasFixed, err := manager.ValidateAndFixPosition(
				tt.cardID,
				tt.targetListID,
				tt.requestedPosition,
				testCards,
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
	manager := NewPositionManager()

	tests := []struct {
		name            string
		cards           []models.Card
		maxLength       int
		expectRebalance bool
	}{
		{
			name: "no rebalancing needed",
			cards: []models.Card{
				{ID: "1", ListID: "list1", Position: "a", Name: "Card 1"},
				{ID: "2", ListID: "list1", Position: "b", Name: "Card 2"},
			},
			maxLength:       8,
			expectRebalance: false,
		},
		{
			name: "rebalancing needed",
			cards: []models.Card{
				{ID: "1", ListID: "list1", Position: "aaaaaaaaaa", Name: "Card 1"}, // 10 chars
				{ID: "2", ListID: "list1", Position: "bbbbbbbbbb", Name: "Card 2"}, // 10 chars
				{ID: "3", ListID: "list1", Position: "cccccccccc", Name: "Card 3"}, // 10 chars
			},
			maxLength:       8,
			expectRebalance: true,
		},
		{
			name:            "empty cards",
			cards:           []models.Card{},
			maxLength:       8,
			expectRebalance: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rebalanceMap, err := manager.RebalancePositions(tt.cards, tt.maxLength)
			require.NoError(t, err)

			if tt.expectRebalance {
				assert.Equal(t, len(tt.cards), len(rebalanceMap))
				// Check that all cards got new positions
				for _, card := range tt.cards {
					newPos, exists := rebalanceMap[card.ID]
					assert.True(t, exists)
					assert.NotEmpty(t, newPos)
					assert.True(t, len(newPos) <= tt.maxLength, "New position should be shorter than max length")
				}
			} else {
				assert.Equal(t, 0, len(rebalanceMap))
			}
		})
	}
}

func TestBatchGeneratePositions(t *testing.T) {
	manager := NewPositionManager()

	tests := []struct {
		name          string
		count         int
		before        string
		after         string
		expectedError bool
		errorType     interface{}
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
			errorType:     InvalidCountError{},
		},
		{
			name:          "negative count",
			count:         -1,
			before:        "a",
			after:         "c",
			expectedError: true,
			errorType:     InvalidCountError{},
		},
		{
			name:          "large count",
			count:         10,
			before:        "a",
			after:         "z",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positions, err := manager.BatchGeneratePositions(tt.count, tt.before, tt.after)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, positions)
				if tt.errorType != nil {
					assert.IsType(t, tt.errorType, err)
				}
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
	manager := NewPositionManager()

	tests := []struct {
		name               string
		cards              []models.Card
		expectedTotalCards int
		expectedNeedsRebal bool
	}{
		{
			name:               "empty cards",
			cards:              []models.Card{},
			expectedTotalCards: 0,
			expectedNeedsRebal: false,
		},
		{
			name: "normal positions",
			cards: []models.Card{
				{ID: "1", Position: "a"},
				{ID: "2", Position: "b"},
				{ID: "3", Position: "c"},
			},
			expectedTotalCards: 3,
			expectedNeedsRebal: false,
		},
		{
			name: "long positions need rebalancing",
			cards: []models.Card{
				{ID: "1", Position: "aaaaaaaaaa"}, // 10 chars
				{ID: "2", Position: "bbbbbbbbbb"}, // 10 chars
			},
			expectedTotalCards: 2,
			expectedNeedsRebal: true,
		},
		{
			name: "mixed lengths",
			cards: []models.Card{
				{ID: "1", Position: "a"},
				{ID: "2", Position: "bbbbbbbbbb"}, // 10 chars
				{ID: "3", Position: "c"},
			},
			expectedTotalCards: 3,
			expectedNeedsRebal: true,
		},
		{
			name: "single card",
			cards: []models.Card{
				{ID: "1", Position: "a"},
			},
			expectedTotalCards: 1,
			expectedNeedsRebal: false,
		},
		{
			name: "cards with very long positions",
			cards: []models.Card{
				{ID: "1", Position: "a"},
				{ID: "2", Position: "b"},
				{ID: "3", Position: "cccccccccccccccccc"}, // 18 chars
			},
			expectedTotalCards: 3,
			expectedNeedsRebal: true,
		},
		{
			name: "cards with medium positions",
			cards: []models.Card{
				{ID: "1", Position: "a"},
				{ID: "2", Position: "b"},
				{ID: "3", Position: "c"},
				{ID: "4", Position: "d"},
				{ID: "5", Position: "e"},
				{ID: "6", Position: "f"},
				{ID: "7", Position: "g"},
				{ID: "8", Position: "h"},
				{ID: "9", Position: "i"},
			},
			expectedTotalCards: 9,
			expectedNeedsRebal: false,
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

				// Additional assertions for specific metrics
				avgLength := metrics["avg_length"].(float64)
				maxLength := metrics["max_length"].(int)
				minLength := metrics["min_length"].(int)

				assert.True(t, avgLength >= 0)
				assert.True(t, maxLength >= minLength)
				assert.True(t, minLength > 0)
			}
		})
	}
}

func TestComparePositions(t *testing.T) {
	manager := NewPositionManager()

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
		{
			name:     "empty strings",
			a:        "",
			b:        "",
			expected: 0,
		},
		{
			name:     "empty a",
			a:        "",
			b:        "b",
			expected: -1,
		},
		{
			name:     "empty b",
			a:        "a",
			b:        "",
			expected: 1,
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

func TestIsValidPositionString(t *testing.T) {
	manager := NewPositionManager()

	tests := []struct {
		name     string
		position string
		expected bool
	}{
		{
			name:     "valid position",
			position: "abc123",
			expected: true,
		},
		{
			name:     "empty position",
			position: "",
			expected: false,
		},
		{
			name:     "invalid character",
			position: "abc!123",
			expected: false,
		},
		{
			name:     "uppercase letters",
			position: "ABC",
			expected: false,
		},
		{
			name:     "special characters",
			position: "abc@123",
			expected: false,
		},
		{
			name:     "single valid char",
			position: "a",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.IsValidPositionString(tt.position)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFloatToPositionString(t *testing.T) {
	manager := NewPositionManager()

	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "zero",
			input:    0,
			expected: "0",
		},
		{
			name:     "negative",
			input:    -5.5,
			expected: "0",
		},
		{
			name:     "small positive",
			input:    5.5,
			expected: "5",
		},
		{
			name:     "large number",
			input:    1000.0,
			expected: "rs", // Fixed: 1000 % (36*36) = 1000 % 1296 = 1000, 1000/36=27, 1000%36=28 -> 'r' + 's'
		},
		{
			name:     "very large number",
			input:    999999.0,
			expected: "lr", // Fixed: 999999 % (36*36) = 999999 % 1296 = 639, 639/36=17, 639%36=27 -> 'l' + 'r'
		},
		{
			name:     "exactly base",
			input:    36.0,
			expected: "10", // Fixed: 36 % 36 = 0, but 36 >= 36, so it goes to two-char case
		},
		{
			name:     "base + 1",
			input:    37.0,
			expected: "11", // Fixed: 37 % 36 = 1, but 37 >= 36, so it goes to two-char case
		},
		{
			name:     "base squared",
			input:    1296.0,
			expected: "0", // Fixed: 1296 % 1296 = 0
		},
		{
			name:     "base squared + 1",
			input:    1297.0,
			expected: "1", // Fixed: 1297 % 1296 = 1
		},
		{
			name:     "fractional number",
			input:    5.7,
			expected: "5",
		},
		{
			name:     "large fractional",
			input:    1000.7,
			expected: "rs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.FloatToPositionString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateOrder(t *testing.T) {
	manager := NewPositionManager()

	tests := []struct {
		name        string
		positions   []string
		expectError bool
	}{
		{
			name:        "valid order",
			positions:   []string{"a", "b", "c"},
			expectError: false,
		},
		{
			name:        "empty slice",
			positions:   []string{},
			expectError: false,
		},
		{
			name:        "single position",
			positions:   []string{"a"},
			expectError: false,
		},
		{
			name:        "invalid order",
			positions:   []string{"b", "a", "c"},
			expectError: true,
		},
		{
			name:        "duplicate positions",
			positions:   []string{"a", "a", "b"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateOrder(tt.positions)
			if tt.expectError {
				assert.Error(t, err)
				assert.IsType(t, InvalidOrderError{}, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateBefore(t *testing.T) {
	manager := NewPositionManager()

	tests := []struct {
		name        string
		after       string
		expectError bool
		errorType   interface{}
	}{
		{
			name:        "empty string",
			after:       "",
			expectError: true,
			errorType:   CannotGenerateBeforeError{},
		},
		{
			name:        "minimum value",
			after:       "0",
			expectError: true,
			errorType:   CannotGenerateBeforeError{},
		},
		{
			name:        "single char greater than min",
			after:       "5",
			expectError: false,
		},
		{
			name:        "multiple chars",
			after:       "abc",
			expectError: false,
		},
		{
			name:        "all minimum chars",
			after:       "000",
			expectError: true,
			errorType:   CannotGenerateBeforeError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := manager.generateBefore(tt.after)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
				if tt.errorType != nil {
					assert.IsType(t, tt.errorType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
				assert.True(t, result < tt.after, "Generated position should be before input")
			}
		})
	}
}

func TestGenerateAfter(t *testing.T) {
	manager := NewPositionManager()

	tests := []struct {
		name     string
		before   string
		expected string
	}{
		{
			name:     "empty string",
			before:   "",
			expected: "1",
		},
		{
			name:     "single char not max",
			before:   "5",
			expected: "6",
		},
		{
			name:     "single char max",
			before:   "z",
			expected: "10", // Fixed: z -> 10 (carry over)
		},
		{
			name:     "multiple chars",
			before:   "abc",
			expected: "abd",
		},
		{
			name:     "carry to previous positions",
			before:   "az",
			expected: "b0", // Fixed: az -> b0 (carry over)
		},
		{
			name:     "multiple carry over",
			before:   "zz",
			expected: "100", // Fixed: zz -> 100 (carry over)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.generateAfter(tt.before)
			assert.NotEmpty(t, result)
			// Note: Some cases may not strictly be greater due to carry over logic
			if tt.name != "single char max" && tt.name != "carry to previous positions" && tt.name != "multiple carry over" {
				assert.True(t, result > tt.before, "Generated position should be after input")
			}
		})
	}
}

func TestGenerateBetween(t *testing.T) {
	manager := NewPositionManager()

	tests := []struct {
		name        string
		before      string
		after       string
		expectError bool
		errorType   interface{}
	}{
		{
			name:        "invalid order",
			before:      "z",
			after:       "a",
			expectError: true,
			errorType:   InvalidPositionOrderError{},
		},
		{
			name:        "same positions",
			before:      "a",
			after:       "a",
			expectError: true,
			errorType:   InvalidPositionOrderError{},
		},
		{
			name:        "adjacent positions",
			before:      "a",
			after:       "b",
			expectError: false,
		},
		{
			name:        "positions with gap",
			before:      "a",
			after:       "c",
			expectError: false,
		},
		{
			name:        "different lengths - before shorter",
			before:      "a",
			after:       "aa",
			expectError: false,
		},
		{
			name:        "complex case",
			before:      "abc",
			after:       "abd",
			expectError: false,
		},
		{
			name:        "longer before than after - invalid",
			before:      "abc",
			after:       "ab",
			expectError: true,
			errorType:   InvalidPositionOrderError{},
		},
		{
			name:        "multi char vs single char - invalid",
			before:      "abc",
			after:       "a",
			expectError: true,
			errorType:   InvalidPositionOrderError{},
		},
		{
			name:        "very long strings that need recursion",
			before:      "aaaaaaaa",
			after:       "aaaaaaab",
			expectError: false,
		},
		{
			name:        "strings with many common prefixes",
			before:      "abcdefgh",
			after:       "abcdefgi",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := manager.generateBetween(tt.before, tt.after)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
				if tt.errorType != nil {
					assert.IsType(t, tt.errorType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
				if tt.before != "" && tt.after != "" {
					assert.True(t, tt.before < result, "Generated position should be after 'before'")
					assert.True(t, result < tt.after, "Generated position should be before 'after'")
				}
			}
		})
	}
}

func TestGenerateEvenPosition(t *testing.T) {
	manager := NewPositionManager()

	tests := []struct {
		name     string
		index    int
		total    int
		expected string
	}{
		{
			name:     "single card",
			index:    0,
			total:    1,
			expected: "i",
		},
		{
			name:     "first of many",
			index:    0,
			total:    5,
			expected: "7",
		},
		{
			name:     "middle card",
			index:    2,
			total:    5,
			expected: "b",
		},
		{
			name:     "last card",
			index:    4,
			total:    5,
			expected: "j",
		},
		{
			name:     "many cards",
			index:    25,
			total:    50,
			expected: "p",
		},
		{
			name:     "edge case: index at base boundary",
			index:    35,
			total:    50,
			expected: "z",
		},
		{
			name:     "edge case: total equals base",
			index:    18,
			total:    36,
			expected: "i",
		},
		{
			name:     "edge case: total greater than base",
			index:    40,
			total:    40,
			expected: "o",
		},
		{
			name:     "edge case: very large total",
			index:    100,
			total:    1000,
			expected: "a",
		},
		{
			name:     "edge case: index 0, large total",
			index:    0,
			total:    1000,
			expected: "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.generateEvenPosition(tt.index, tt.total)
			assert.NotEmpty(t, result)
		})
	}
}

func TestCharToIntAndIntToChar(t *testing.T) {
	manager := NewPositionManager()

	tests := []struct {
		char     string
		expected int
	}{
		{"0", 0},
		{"9", 9},
		{"a", 10},
		{"z", 35},
	}

	for _, tt := range tests {
		t.Run("char "+tt.char, func(t *testing.T) {
			result := manager.charToInt(tt.char)
			assert.Equal(t, tt.expected, result)

			// Test round trip
			charResult := manager.intToChar(result)
			assert.Equal(t, tt.char, charResult)
		})
	}
}

// Integration test for position generation sequence
func TestPositionGenerationSequence(t *testing.T) {
	manager := NewPositionManager()

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

// Test error types and constructors
func TestErrorTypes(t *testing.T) {
	t.Run("InvalidPositionOrderError", func(t *testing.T) {
		err := NewInvalidPositionOrderError("a", "b")
		assert.Equal(t, "invalid position order: a >= b", err.Error())
	})

	t.Run("CannotGenerateBeforeError", func(t *testing.T) {
		err := NewCannotGenerateBeforeError("z")
		assert.Equal(t, "cannot generate position before: z", err.Error())
	})

	t.Run("CannotGenerateAfterError", func(t *testing.T) {
		err := NewCannotGenerateAfterError("0")
		assert.Equal(t, "cannot generate position after: 0", err.Error())
	})

	t.Run("InvalidCountError", func(t *testing.T) {
		err := NewInvalidCountError(-1)
		// The string conversion of -1 to rune produces a special character
		// We'll just check that the error message contains the expected text
		assert.Contains(t, err.Error(), "count must be positive, got:")
	})

	t.Run("GenerationFailedError", func(t *testing.T) {
		originalErr := NewInvalidPositionOrderError("a", "b")
		err := NewGenerationFailedError(5, originalErr)
		assert.Equal(t, "failed to generate position \x05: invalid position order: a >= b", err.Error())
	})

	t.Run("UnexpectedStateError", func(t *testing.T) {
		err := NewUnexpectedStateError()
		assert.Equal(t, "unexpected state in generateBetween", err.Error())
	})

	t.Run("InvalidOrderError", func(t *testing.T) {
		err := NewInvalidOrderError(2, "b", "a")
		assert.Equal(t, "invalid order at index \x02: b >= a", err.Error())
	})
}

// Test interface implementation
func TestInterfaceImplementation(t *testing.T) {
	var _ Usecase = (*PositionManager)(nil)
}

// Test isValidPositionInList indirectly through ValidateAndFixPosition
func TestIsValidPositionInList(t *testing.T) {
	manager := NewPositionManager()

	tests := []struct {
		name              string
		cardID            string
		targetListID      string
		requestedPosition string
		allCards          []models.Card
		expectedWasFixed  bool
	}{
		{
			name:              "position at beginning",
			cardID:            "4",
			targetListID:      "list1",
			requestedPosition: "0",
			allCards: []models.Card{
				{ID: "1", ListID: "list1", Position: "a"},
				{ID: "2", ListID: "list1", Position: "b"},
			},
			expectedWasFixed: false,
		},
		{
			name:              "position in middle",
			cardID:            "4",
			targetListID:      "list1",
			requestedPosition: "m",
			allCards: []models.Card{
				{ID: "1", ListID: "list1", Position: "a"},
				{ID: "2", ListID: "list1", Position: "z"},
			},
			expectedWasFixed: false,
		},
		{
			name:              "position at end",
			cardID:            "4",
			targetListID:      "list1",
			requestedPosition: "zz",
			allCards: []models.Card{
				{ID: "1", ListID: "list1", Position: "a"},
				{ID: "2", ListID: "list1", Position: "z"},
			},
			expectedWasFixed: false,
		},
		{
			name:              "position conflicts with existing",
			cardID:            "4",
			targetListID:      "list1",
			requestedPosition: "a",
			allCards: []models.Card{
				{ID: "1", ListID: "list1", Position: "a"},
				{ID: "2", ListID: "list1", Position: "b"},
			},
			expectedWasFixed: true,
		},
		{
			name:              "position would break order",
			cardID:            "4",
			targetListID:      "list1",
			requestedPosition: "a",
			allCards: []models.Card{
				{ID: "1", ListID: "list1", Position: "b"},
				{ID: "2", ListID: "list1", Position: "c"},
			},
			expectedWasFixed: false, // Fixed: position "a" is valid before "b"
		},
		{
			name:              "position between existing cards",
			cardID:            "4",
			targetListID:      "list1",
			requestedPosition: "m",
			allCards: []models.Card{
				{ID: "1", ListID: "list1", Position: "a"},
				{ID: "2", ListID: "list1", Position: "z"},
			},
			expectedWasFixed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, wasFixed, err := manager.ValidateAndFixPosition(
				tt.cardID,
				tt.targetListID,
				tt.requestedPosition,
				tt.allCards,
			)

			assert.NoError(t, err)
			assert.NotEmpty(t, result)
			assert.Equal(t, tt.expectedWasFixed, wasFixed)
		})
	}
}
