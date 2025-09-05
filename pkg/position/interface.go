package position

import "github.com/nguyentantai21042004/kanban-api/internal/models"

//go:generate mockery --name UseCase
type Usecase interface {
	// GeneratePosition creates a position between two existing positions
	GeneratePosition(before, after string) (string, error)

	// ValidateAndFixPosition validates a requested position and fixes if needed
	ValidateAndFixPosition(cardID, targetListID, requestedPosition string, allCards []models.Card) (string, bool, error)

	// RebalancePositions generates new positions for cards when they become too long
	RebalancePositions(cards []models.Card, maxLength int) (map[string]string, error)

	// BatchGeneratePositions generates positions for multiple cards efficiently
	BatchGeneratePositions(count int, before, after string) ([]string, error)

	// ComparePositions compares two position strings (-1, 0, 1)
	ComparePositions(a, b string) int

	// IsValidPositionString is an exported validator for external callers
	IsValidPositionString(position string) bool

	// FloatToPositionString converts a legacy numeric/float position into a base-36 string
	FloatToPositionString(floatPos float64) string

	// ValidateOrder validates that a slice of positions is in correct order
	ValidateOrder(positions []string) error

	// GetPositionMetrics returns metrics about position distribution
	GetPositionMetrics(cards []models.Card) map[string]interface{}
}

// Custom error types for position management
type (
	// InvalidPositionOrderError occurs when trying to generate a position between invalid order
	InvalidPositionOrderError struct {
		Before string
		After  string
	}

	// CannotGenerateBeforeError occurs when trying to generate a position before the minimum
	CannotGenerateBeforeError struct {
		Position string
	}

	// CannotGenerateAfterError occurs when trying to generate a position after the maximum
	CannotGenerateAfterError struct {
		Position string
	}

	// InvalidCountError occurs when count parameter is invalid
	InvalidCountError struct {
		Count int
	}

	// GenerationFailedError occurs when position generation fails
	GenerationFailedError struct {
		Index int
		Err   error
	}

	// UnexpectedStateError occurs when an unexpected state is reached during generation
	UnexpectedStateError struct{}

	// InvalidOrderError occurs when positions are not in correct order
	InvalidOrderError struct {
		Index  int
		Before string
		After  string
	}
)

// Error implementations
func (e InvalidPositionOrderError) Error() string {
	return "invalid position order: " + e.Before + " >= " + e.After
}

func (e CannotGenerateBeforeError) Error() string {
	return "cannot generate position before: " + e.Position
}

func (e CannotGenerateAfterError) Error() string {
	return "cannot generate position after: " + e.Position
}

func (e InvalidCountError) Error() string {
	return "count must be positive, got: " + string(rune(e.Count))
}

func (e GenerationFailedError) Error() string {
	return "failed to generate position " + string(rune(e.Index)) + ": " + e.Err.Error()
}

func (e UnexpectedStateError) Error() string {
	return "unexpected state in generateBetween"
}

func (e InvalidOrderError) Error() string {
	return "invalid order at index " + string(rune(e.Index)) + ": " + e.Before + " >= " + e.After
}

// Error constructors
func NewInvalidPositionOrderError(before, after string) InvalidPositionOrderError {
	return InvalidPositionOrderError{Before: before, After: after}
}

func NewCannotGenerateBeforeError(position string) CannotGenerateBeforeError {
	return CannotGenerateBeforeError{Position: position}
}

func NewCannotGenerateAfterError(position string) CannotGenerateAfterError {
	return CannotGenerateAfterError{Position: position}
}

func NewInvalidCountError(count int) InvalidCountError {
	return InvalidCountError{Count: count}
}

func NewGenerationFailedError(index int, err error) GenerationFailedError {
	return GenerationFailedError{Index: index, Err: err}
}

func NewUnexpectedStateError() UnexpectedStateError {
	return UnexpectedStateError{}
}

func NewInvalidOrderError(index int, before, after string) InvalidOrderError {
	return InvalidOrderError{Index: index, Before: before, After: after}
}
