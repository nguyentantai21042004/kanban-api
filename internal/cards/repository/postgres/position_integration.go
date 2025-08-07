package postgres

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
	"gitlab.com/tantai-kanban/kanban-api/internal/cards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/position"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

// EnhancedMove implements advanced position management for card moves
func (r implRepository) EnhancedMove(ctx context.Context, sc models.Scope, opts repository.MoveOptions) (models.Card, error) {
	// Start transaction
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.EnhancedMove.BeginTx: %v", err)
		return models.Card{}, err
	}
	defer tx.Rollback()

	// Initialize position manager
	positionManager := position.NewManager()

	// Get all cards in target list (excluding the card being moved)
	allCards, err := r.getCardsInListExcluding(ctx, opts.ListID, opts.ID)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.EnhancedMove.getCardsInList: %v", err)
		return models.Card{}, err
	}

	// Convert to position manager format
	posCards := make([]position.Card, len(allCards))
	for i, card := range allCards {
		positionStr := r.extractPositionString(card)
		posCards[i] = position.Card{
			ID:       card.ID,
			ListID:   card.ListID,
			Position: positionStr,
			Name:     card.Name,
		}
	}

	// Validate and potentially fix the requested position
	requestedPosition := r.extractPositionFromOptions(opts)
	finalPositionStr, wasFixed, err := positionManager.ValidateAndFixPosition(
		opts.ID,
		opts.ListID,
		requestedPosition,
		posCards,
	)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.EnhancedMove.ValidateAndFixPosition: %v", err)
		return models.Card{}, err
	}

	if wasFixed {
		r.l.Infof(ctx, "Position was optimized by server: %s -> %s", requestedPosition, finalPositionStr)
	}

	// Convert position string back to decimal for storage
	finalPosition := r.convertPositionStringToDecimal(finalPositionStr)

	// Check if rebalancing is needed
	metrics := positionManager.GetPositionMetrics(posCards)
	if needsRebalance, ok := metrics["needs_rebalance"].(bool); ok && needsRebalance {
		r.l.Infof(ctx, "Position rebalancing recommended for list %s", opts.ListID)
		// Note: Rebalancing should be done in a background job to avoid blocking the move
	}

	// Convert final position back to float64 for MoveOptions
	finalPositionFloat, _ := finalPosition.Big.Float64()

	// Build the move model with the final position
	c, col, err := r.buildEnhancedMoveModel(ctx, repository.MoveOptions{
		ID:       opts.ID,
		ListID:   opts.ListID,
		Position: finalPositionFloat,
		OldModel: opts.OldModel,
	})
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.EnhancedMove.buildMoveModel: %v", err)
		return models.Card{}, err
	}

	// Update the card
	_, err = c.Update(ctx, tx, boil.Whitelist(col...))
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.EnhancedMove.Update: %v", err)
		return models.Card{}, err
	}

	// Create enhanced activity record
	activity := r.buildEnhancedActivityModel(ctx, c.ID, string(models.CardActionTypeMoved),
		map[string]interface{}{
			"from_list_id":  opts.OldModel.ListID,
			"from_position": r.extractPositionString(r.convertModelToDBModel(opts.OldModel)),
		},
		map[string]interface{}{
			"to_list_id":    opts.ListID,
			"to_position":   finalPositionStr,
			"was_optimized": wasFixed,
		})

	err = activity.Insert(ctx, tx, boil.Infer())
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.EnhancedMove.InsertActivity: %v", err)
		return models.Card{}, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.EnhancedMove.Commit: %v", err)
		return models.Card{}, err
	}

	r.l.Infof(ctx, "Enhanced move completed: %s to %s at position %s", opts.ID, opts.ListID, finalPositionStr)
	return models.NewCard(c), nil
}

// Helper functions for position management

func (r implRepository) getCardsInListExcluding(ctx context.Context, listID, excludeCardID string) ([]dbmodels.Card, error) {
	cards, err := dbmodels.Cards(
		dbmodels.CardWhere.ListID.EQ(listID),
		dbmodels.CardWhere.ID.NEQ(excludeCardID),
	).All(ctx, r.database)
	if err != nil {
		return nil, err
	}

	// Convert to slice of values
	result := make([]dbmodels.Card, len(cards))
	for i, card := range cards {
		result[i] = *card
	}
	return result, nil
}

func (r implRepository) extractPositionString(card dbmodels.Card) string {
	if card.Position.Big == nil {
		return "n" // Default middle position
	}

	// Convert decimal to position string
	floatPos, _ := card.Position.Big.Float64()

	// If it's a simple number, convert to position string
	if floatPos == float64(int64(floatPos)) {
		// Integer position, convert to fractional indexing
		return r.convertFloatToPositionString(floatPos)
	}

	// Already a position string or complex decimal
	return card.Position.Big.String()
}

func (r implRepository) extractPositionFromOptions(opts repository.MoveOptions) string {
	if opts.Position == 0 {
		return ""
	}
	return r.convertFloatToPositionString(opts.Position)
}

func (r implRepository) convertPositionStringToDecimal(posStr string) types.Decimal {
	// Try to parse as float first
	if floatVal, err := strconv.ParseFloat(posStr, 64); err == nil {
		return types.Decimal{Big: decimal.New(int64(floatVal*100), -2)}
	}

	// If it's a fractional indexing string, convert to float representation
	floatVal := r.convertPositionStringToFloat(posStr)
	return types.Decimal{Big: decimal.New(int64(floatVal*100), -2)}
}

func (r implRepository) convertFloatToPositionString(floatPos float64) string {
	// Simple conversion: use base62 encoding for fractional indexing
	base62 := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if floatPos <= 0 {
		return "0"
	}

	// Convert to base62-like representation
	intPos := int(floatPos) % (len(base62) * len(base62)) // Prevent overflow
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

func (r implRepository) convertPositionStringToFloat(posStr string) float64 {
	base62 := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if posStr == "" {
		return 1000.0 // Default position
	}

	// Handle single character positions
	if len(posStr) == 1 {
		index := strings.Index(base62, string(posStr[0]))
		if index >= 0 {
			return float64(index * 100) // Scale for storage
		}
	}

	// Handle multi-character positions
	result := 0.0
	for i, char := range posStr {
		index := strings.Index(base62, string(char))
		if index >= 0 {
			result += float64(index) * pow(62, len(posStr)-i-1)
		}
	}

	return result * 100 // Scale for storage
}

func (r implRepository) buildEnhancedMoveModel(ctx context.Context, opts repository.MoveOptions) (dbmodels.Card, []string, error) {
	card := dbmodels.Card{
		ListID:   opts.ListID,
		Position: types.Decimal{Big: decimal.New(int64(opts.Position*100), -2)},
	}

	cols := []string{
		dbmodels.CardColumns.ListID,
		dbmodels.CardColumns.Position,
		dbmodels.CardColumns.UpdatedAt,
	}

	if err := postgres.IsUUID(opts.ID); err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.buildEnhancedMoveModel.IsUUID: %v", err)
		return dbmodels.Card{}, nil, err
	}
	card.ID = opts.ID
	card.UpdatedAt = r.clock()

	return card, cols, nil
}

func (r implRepository) buildEnhancedActivityModel(ctx context.Context, cardID string, actionType string, oldData, newData map[string]interface{}) dbmodels.CardActivity {
	activity := dbmodels.CardActivity{
		CardID:     cardID,
		ActionType: dbmodels.CardActionType(actionType),
	}

	if oldData != nil {
		oldDataJSON, _ := json.Marshal(oldData)
		activity.OldData = null.JSONFrom(oldDataJSON)
	}

	if newData != nil {
		newDataJSON, _ := json.Marshal(newData)
		activity.NewData = null.JSONFrom(newDataJSON)
	}

	return activity
}

func (r implRepository) convertModelToDBModel(model models.Card) dbmodels.Card {
	return dbmodels.Card{
		ID:       model.ID,
		ListID:   model.ListID,
		Position: types.Decimal{Big: decimal.New(int64(model.Position*100), -2)},
		Name:     model.Name,
	}
}

// Helper function for power calculation
func pow(base, exp int) float64 {
	result := 1.0
	for i := 0; i < exp; i++ {
		result *= float64(base)
	}
	return result
}

// RebalanceListPositions rebalances all positions in a list
func (r implRepository) RebalanceListPositions(ctx context.Context, listID string) error {
	positionManager := position.NewManager()

	// Get all cards in list
	allCards, err := r.getCardsInListExcluding(ctx, listID, "")
	if err != nil {
		return err
	}

	// Convert to position format
	posCards := make([]position.Card, len(allCards))
	for i, card := range allCards {
		posCards[i] = position.Card{
			ID:       card.ID,
			ListID:   card.ListID,
			Position: r.extractPositionString(card),
			Name:     card.Name,
		}
	}

	// Generate rebalanced positions
	rebalanceMap, err := positionManager.RebalancePositions(posCards, 8)
	if err != nil {
		return err
	}

	if len(rebalanceMap) == 0 {
		return nil // No rebalancing needed
	}

	// Start transaction for batch update
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update positions
	for cardID, newPositionStr := range rebalanceMap {
		newPosition := r.convertPositionStringToDecimal(newPositionStr)

		_, err = dbmodels.Cards(dbmodels.CardWhere.ID.EQ(cardID)).UpdateAll(ctx, tx, dbmodels.M{
			dbmodels.CardColumns.Position:  newPosition,
			dbmodels.CardColumns.UpdatedAt: r.clock(),
		})
		if err != nil {
			r.l.Errorf(ctx, "Failed to update card position during rebalance: %v", err)
			return err
		}
	}

	return tx.Commit()
}
