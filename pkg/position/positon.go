package position

import (
	"math"
	"sort"
	"strings"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

const (
	// Character set for base-36 encoding
	BaseChars = "0123456789abcdefghijklmnopqrstuvwxyz"
	Base      = 36
)

// PositionManager handles position calculations using base-36 fractional indexing
type PositionManager struct {
	baseChars string
	base      int
	midChar   string
	minChar   string
	maxChar   string
}

// Ensure PositionManager implements Usecase
var _ Usecase = (*PositionManager)(nil)

// NewPositionManager creates a new position manager
func NewPositionManager() *PositionManager {
	return &PositionManager{
		baseChars: BaseChars,
		base:      Base,
		midChar:   "i", // Middle character in base-36
		minChar:   "0",
		maxChar:   "z",
	}
}

// GeneratePosition creates a position between two existing positions
func (m *PositionManager) GeneratePosition(before, after string) (string, error) {
	// Case 1: Both empty - generate initial key
	if before == "" && after == "" {
		return m.midChar, nil
	}

	// Case 2: Only before is empty - generate key before 'after'
	if before == "" {
		return m.generateBefore(after)
	}

	// Case 3: Only after is empty - generate key after 'before'
	if after == "" {
		return m.generateAfter(before), nil
	}

	// Case 4: Both have values - generate key between them
	if before >= after {
		return "", NewInvalidPositionOrderError(before, after)
	}

	return m.generateBetween(before, after)
}

// ValidateAndFixPosition validates a requested position and fixes if needed
func (m *PositionManager) ValidateAndFixPosition(
	cardID string,
	targetListID string,
	requestedPosition string,
	allCards []models.Card,
) (string, bool, error) {
	// Filter cards in target list (excluding the card being moved)
	listCards := make([]models.Card, 0)
	for _, card := range allCards {
		if card.ListID == targetListID && card.ID != cardID {
			listCards = append(listCards, card)
		}
	}

	// Sort by position
	sort.Slice(listCards, func(i, j int) bool {
		return strings.Compare(listCards[i].Position, listCards[j].Position) < 0
	})

	// If no requested position or invalid, generate optimal position
	if requestedPosition == "" || !m.isValidPosition(requestedPosition) {
		optimalPosition := m.generateOptimalPosition(listCards)
		return optimalPosition, true, nil
	}

	// Check if requested position maintains sort order
	if m.isValidPositionInList(requestedPosition, listCards) {
		// Additional check: if position is the same as an existing card, consider it invalid
		for _, card := range listCards {
			if card.Position == requestedPosition {
				// Position conflicts with existing card, generate better one
				optimalPosition := m.generateOptimalPosition(listCards)
				return optimalPosition, true, nil
			}
		}

		// For empty list, always generate optimal position to ensure consistency
		if len(listCards) == 0 {
			optimalPosition := m.generateOptimalPosition(listCards)
			return optimalPosition, true, nil
		}

		return requestedPosition, false, nil
	}

	// Position would break order, generate better one
	optimalPosition := m.generateOptimalPosition(listCards)
	return optimalPosition, true, nil
}

// RebalancePositions generates new positions for cards when they become too long
func (m *PositionManager) RebalancePositions(cards []models.Card, maxLength int) (map[string]string, error) {
	rebalanceMap := make(map[string]string)

	// Check if rebalancing is needed
	needsRebalancing := false
	for _, card := range cards {
		if len(card.Position) > maxLength {
			needsRebalancing = true
			break
		}
	}

	if !needsRebalancing {
		return rebalanceMap, nil
	}

	// Sort cards by current position
	sortedCards := make([]models.Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return strings.Compare(sortedCards[i].Position, sortedCards[j].Position) < 0
	})

	// Generate new evenly spaced positions
	for i, card := range sortedCards {
		newPosition := m.generateEvenPosition(i, len(sortedCards))
		rebalanceMap[card.ID] = newPosition
	}

	return rebalanceMap, nil
}

// BatchGeneratePositions generates positions for multiple cards efficiently
func (m *PositionManager) BatchGeneratePositions(count int, before, after string) ([]string, error) {
	if count <= 0 {
		return nil, NewInvalidCountError(count)
	}

	if count == 1 {
		pos, err := m.GeneratePosition(before, after)
		if err != nil {
			return nil, err
		}
		return []string{pos}, nil
	}

	positions := make([]string, count)
	currentBefore := before

	for i := 0; i < count; i++ {
		currentAfter := ""
		if i == count-1 {
			currentAfter = after
		}

		pos, err := m.GeneratePosition(currentBefore, currentAfter)
		if err != nil {
			return nil, NewGenerationFailedError(i, err)
		}

		positions[i] = pos
		currentBefore = pos
	}

	return positions, nil
}

// ComparePositions compares two position strings (-1, 0, 1)
func (m *PositionManager) ComparePositions(a, b string) int {
	return strings.Compare(a, b)
}

// IsValidPosition checks if a position string is valid
func (m *PositionManager) isValidPosition(position string) bool {
	if position == "" {
		return false
	}

	for _, char := range position {
		if !strings.ContainsRune(m.baseChars, char) {
			return false
		}
	}

	return true
}

// IsValidPositionString is an exported validator for external callers (e.g. HTTP layer)
func (m *PositionManager) IsValidPositionString(position string) bool {
	return m.isValidPosition(position)
}

// FloatToPositionString converts a legacy numeric/float position into a base-36 string.
// This is used during migration/backward-compat scenarios to provide a stable string representation
// for API responses and websocket payloads when only a float value is available.
func (m *PositionManager) FloatToPositionString(floatPos float64) string {
	if floatPos <= 0 {
		return "0"
	}

	// Convert to a base-36 representation
	intPos := int(floatPos) % (Base * Base)
	if intPos >= Base {
		// Use two characters for larger positions
		first := intPos / Base
		second := intPos % Base
		if first < Base && second < Base {
			return string(BaseChars[first]) + string(BaseChars[second])
		}
	}

	if intPos < Base {
		return string(BaseChars[intPos])
	}

	// Fallback for very large numbers
	return "z"
}

// isValidPositionInList checks if position maintains order in the list
func (m *PositionManager) isValidPositionInList(position string, existingCards []models.Card) bool {
	if len(existingCards) == 0 {
		return true
	}

	// Find the correct insertion point
	for i, card := range existingCards {
		if strings.Compare(position, card.Position) <= 0 {
			// Position should be inserted here
			if i > 0 && strings.Compare(existingCards[i-1].Position, position) >= 0 {
				return false // Would break order
			}
			return true
		}
	}

	// Position would be inserted at the end
	// Check if it's greater than the last card
	if len(existingCards) > 0 {
		return strings.Compare(position, existingCards[len(existingCards)-1].Position) > 0
	}

	return true
}

// generateOptimalPosition creates an optimal position for the end of the list
func (m *PositionManager) generateOptimalPosition(existingCards []models.Card) string {
	if len(existingCards) == 0 {
		return m.midChar
	}

	// Append to end
	lastCard := existingCards[len(existingCards)-1]
	position := m.generateAfter(lastCard.Position)
	return position
}

// generateBefore creates a position before the given position
func (m *PositionManager) generateBefore(after string) (string, error) {
	if after == "" || after == "0" {
		return "", NewCannotGenerateBeforeError(after)
	}

	// Strategy: Decrement last character, append max char
	for i := len(after) - 1; i >= 0; i-- {
		char := string(after[i])
		charInt := m.charToInt(char)

		if charInt > 0 {
			newChar := m.intToChar(charInt - 1)
			return after[:i] + newChar + "z", nil
		}
	}

	return "", NewCannotGenerateBeforeError(after)
}

// generateAfter creates a position after the given position
func (m *PositionManager) generateAfter(before string) string {
	if before == "" {
		return "1"
	}

	lastChar := string(before[len(before)-1])
	lastInt := m.charToInt(lastChar)

	if lastInt < m.base-1 {
		return before[:len(before)-1] + m.intToChar(lastInt+1)
	}

	// Carry to previous positions
	return m.generateAfter(before[:len(before)-1]) + "0"
}

// generateBetween creates a position between two positions using base-36 logic
func (m *PositionManager) generateBetween(before, after string) (string, error) {
	if before >= after {
		return "", NewInvalidPositionOrderError(before, after)
	}

	result := ""
	maxLen := int(math.Max(float64(len(before)), float64(len(after))))

	for i := 0; i < maxLen; i++ {
		aChar := "0"
		bChar := "0"

		if i < len(before) {
			aChar = string(before[i])
		}
		if i < len(after) {
			bChar = string(after[i])
		}

		aInt := m.charToInt(aChar)
		bInt := m.charToInt(bChar)

		if aInt == bInt {
			result += aChar
			continue
		}

		if bInt-aInt > 1 {
			// Found gap, can insert between
			midInt := (aInt + bInt) / 2
			result += m.intToChar(midInt)
			return result, nil
		}

		// bInt - aInt == 1, need to extend string
		result += aChar

		// Handle suffixes
		aSuffix := ""
		bSuffix := ""

		if i+1 < len(before) {
			aSuffix = before[i+1:]
		}
		if i+1 < len(after) {
			bSuffix = after[i+1:]
		}

		if aSuffix == "" && bSuffix == "" {
			result += m.midChar
			return result, nil
		}

		if aSuffix == "" {
			before, err := m.generateBefore(bSuffix)
			if err != nil {
				return "", err
			}
			result += before
			return result, nil
		}

		if bSuffix == "" {
			result += m.generateAfter(aSuffix)
			return result, nil
		}

		between, err := m.generateBetween(aSuffix, bSuffix)
		if err != nil {
			return "", err
		}
		result += between
		return result, nil
	}

	return "", NewUnexpectedStateError()
}

// generateEvenPosition creates evenly spaced positions for rebalancing
func (m *PositionManager) generateEvenPosition(index, total int) string {
	if total == 1 {
		return m.midChar
	}

	// Create evenly spaced positions using base-36
	spacing := m.base / (total + 1)
	charIndex := (index + 1) * spacing

	if charIndex >= m.base {
		charIndex = m.base - 1
	}

	char := m.intToChar(charIndex)

	// For more cards, use two characters
	if total > m.base {
		secondIndex := index % m.base
		secondChar := m.intToChar(secondIndex)
		return char + secondChar
	}

	return char
}

// GetPositionMetrics returns metrics about position distribution
func (m *PositionManager) GetPositionMetrics(cards []models.Card) map[string]interface{} {
	if len(cards) == 0 {
		return map[string]interface{}{
			"total_cards":     0,
			"avg_length":      0,
			"max_length":      0,
			"min_length":      0,
			"needs_rebalance": false,
		}
	}

	totalLength := 0
	maxLength := 0
	minLength := len(cards[0].Position)

	for _, card := range cards {
		length := len(card.Position)
		totalLength += length

		if length > maxLength {
			maxLength = length
		}
		if length < minLength {
			minLength = length
		}
	}

	avgLength := float64(totalLength) / float64(len(cards))
	needsRebalance := maxLength > 8 || avgLength > 5

	return map[string]interface{}{
		"total_cards":     len(cards),
		"avg_length":      avgLength,
		"max_length":      maxLength,
		"min_length":      minLength,
		"needs_rebalance": needsRebalance,
	}
}

// Helper methods for base-36 conversion
func (m *PositionManager) charToInt(char string) int {
	return strings.Index(m.baseChars, char)
}

func (m *PositionManager) intToChar(i int) string {
	return string(m.baseChars[i])
}

// ValidateOrder validates that a slice of positions is in correct order
func (m *PositionManager) ValidateOrder(positions []string) error {
	for i := 1; i < len(positions); i++ {
		if positions[i-1] >= positions[i] {
			return NewInvalidOrderError(i, positions[i-1], positions[i])
		}
	}
	return nil
}
