package position

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Manager handles position calculations using fractional indexing
type Manager struct {
	base62Chars string
	midPoint    string
	minChar     string
	maxChar     string
}

// Card represents a positionable item
type Card struct {
	ID       string
	ListID   string
	Position string
	Name     string
}

// NewManager creates a new position manager
func NewManager() *Manager {
	return &Manager{
		base62Chars: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		midPoint:    "n",
		minChar:     "0",
		maxChar:     "z",
	}
}

// GeneratePosition creates a position between two existing positions
func (m *Manager) GeneratePosition(before, after string) (string, error) {
	if before == "" && after == "" {
		return m.midPoint, nil
	}

	if before == "" {
		return m.generateBefore(after)
	}

	if after == "" {
		return m.generateAfter(before)
	}

	if before >= after {
		return "", fmt.Errorf("invalid position order: %s >= %s", before, after)
	}

	return m.generateBetween(before, after)
}

// ValidateAndFixPosition validates a requested position and fixes if needed
func (m *Manager) ValidateAndFixPosition(
	cardID string,
	targetListID string,
	requestedPosition string,
	allCards []Card,
) (string, bool, error) {
	// Filter cards in target list (excluding the card being moved)
	listCards := make([]Card, 0)
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
		return requestedPosition, false, nil
	}

	// Position would break order, generate better one
	optimalPosition := m.generateOptimalPosition(listCards)
	return optimalPosition, true, nil
}

// RebalancePositions generates new positions for cards when they become too long
func (m *Manager) RebalancePositions(cards []Card, maxLength int) (map[string]string, error) {
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
	sortedCards := make([]Card, len(cards))
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
func (m *Manager) BatchGeneratePositions(count int, before, after string) ([]string, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be positive")
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
			return nil, fmt.Errorf("failed to generate position %d: %w", i, err)
		}

		positions[i] = pos
		currentBefore = pos
	}

	return positions, nil
}

// ComparePositions compares two position strings (-1, 0, 1)
func (m *Manager) ComparePositions(a, b string) int {
	return strings.Compare(a, b)
}

// IsValidPosition checks if a position string is valid
func (m *Manager) isValidPosition(position string) bool {
	if position == "" {
		return false
	}

	for _, char := range position {
		if !strings.ContainsRune(m.base62Chars, char) {
			return false
		}
	}

	return true
}

// IsValidPositionString is an exported validator for external callers (e.g. HTTP layer)
func (m *Manager) IsValidPositionString(position string) bool {
    return m.isValidPosition(position)
}

// FloatToPositionString converts a legacy numeric/float position into a fractional base62 string.
// This is used during migration/backward-compat scenarios to provide a stable string representation
// for API responses and websocket payloads when only a float value is available.
func FloatToPositionString(floatPos float64) string {
    base62 := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

    if floatPos <= 0 {
        return "0"
    }

    // Convert to a base62-like representation. This provides a deterministic ordering-preserving
    // mapping for typical integer-like positions used previously.
    intPos := int(floatPos) % (len(base62) * len(base62))
    if intPos >= len(base62) {
        // Use two characters for larger positions
        first := intPos / len(base62)
        second := intPos % len(base62)
        if first < len(base62) && second < len(base62) {
            return string(base62[first]) + string(base62[second])
        }
    }

    if intPos < len(base62) {
        return string(base62[intPos])
    }

    // Fallback for very large numbers
    return "z"
}

// isValidPositionInList checks if position maintains order in the list
func (m *Manager) isValidPositionInList(position string, existingCards []Card) bool {
	for i, card := range existingCards {
		if strings.Compare(position, card.Position) <= 0 {
			// Position should be inserted here
			if i > 0 && strings.Compare(existingCards[i-1].Position, position) >= 0 {
				return false // Would break order
			}
			return true
		}
	}
	return true // Position at end is valid
}

// generateOptimalPosition creates an optimal position for the end of the list
func (m *Manager) generateOptimalPosition(existingCards []Card) string {
	if len(existingCards) == 0 {
		return m.midPoint
	}

	// Append to end
	lastCard := existingCards[len(existingCards)-1]
	position, _ := m.generateAfter(lastCard.Position)
	return position
}

// generateBefore creates a position before the given position
func (m *Manager) generateBefore(after string) (string, error) {
	if after == "" {
		return m.minChar, nil
	}

	firstChar := string(after[0])
	index := strings.Index(m.base62Chars, firstChar)

	if index > 0 {
		// Can decrement first character
		return string(m.base62Chars[index-1]) + after[1:], nil
	}

	// First character is minimum, need to prepend
	before, err := m.generateBefore(after)
	if err != nil {
		return "", err
	}
	return m.minChar + before, nil
}

// generateAfter creates a position after the given position
func (m *Manager) generateAfter(before string) (string, error) {
	if before == "" {
		return m.midPoint, nil
	}

	// Try to increment last character
	lastChar := string(before[len(before)-1])
	index := strings.Index(m.base62Chars, lastChar)

	if index < len(m.base62Chars)-1 {
		// Can increment last character
		return before[:len(before)-1] + string(m.base62Chars[index+1]), nil
	}

	// Last character is maximum, need to append
	return before + m.midPoint, nil
}

// generateBetween creates a position between two positions
func (m *Manager) generateBetween(before, after string) (string, error) {
	if before >= after {
		return "", fmt.Errorf("invalid position order: %s >= %s", before, after)
	}

	// Find first differing position
	i := 0
	minLength := int(math.Min(float64(len(before)), float64(len(after))))

	for i < minLength && before[i] == after[i] {
		i++
	}

	commonPrefix := before[:i]
	var beforeSuffix, afterSuffix string

	if i < len(before) {
		beforeSuffix = before[i:]
	} else {
		beforeSuffix = m.minChar
	}

	if i < len(after) {
		afterSuffix = after[i:]
	} else {
		afterSuffix = m.maxChar
	}

	beforeChar := string(beforeSuffix[0])
	afterChar := string(afterSuffix[0])

	beforeIndex := strings.Index(m.base62Chars, beforeChar)
	afterIndex := strings.Index(m.base62Chars, afterChar)

	if afterIndex-beforeIndex > 1 {
		// Can insert character between
		midIndex := (beforeIndex + afterIndex) / 2
		return commonPrefix + string(m.base62Chars[midIndex]), nil
	}

	if len(beforeSuffix) > 1 {
		// Extend before suffix
		after, err := m.generateAfter(beforeSuffix[1:])
		if err != nil {
			return "", err
		}
		return commonPrefix + beforeChar + after, nil
	}

	if len(afterSuffix) > 1 {
		// Use after suffix
		before, err := m.generateBefore(afterSuffix[1:])
		if err != nil {
			return "", err
		}
		return commonPrefix + beforeChar + before, nil
	}

	// Both are single characters and adjacent, need to go deeper
	return commonPrefix + beforeChar + m.midPoint, nil
}

// generateEvenPosition creates evenly spaced positions for rebalancing
func (m *Manager) generateEvenPosition(index, total int) string {
	if total == 1 {
		return m.midPoint
	}

	// Create evenly spaced positions using base62
	spacing := len(m.base62Chars) / (total + 1)
	charIndex := (index + 1) * spacing

	if charIndex >= len(m.base62Chars) {
		charIndex = len(m.base62Chars) - 1
	}

	char := string(m.base62Chars[charIndex])

	// For more cards, use two characters
	if total > len(m.base62Chars) {
		secondIndex := index % len(m.base62Chars)
		secondChar := string(m.base62Chars[secondIndex])
		return char + secondChar
	}

	return char
}

// GetPositionMetrics returns metrics about position distribution
func (m *Manager) GetPositionMetrics(cards []Card) map[string]interface{} {
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
