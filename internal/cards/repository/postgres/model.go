package postgres

import (
	"context"
	"encoding/json"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
	"gitlab.com/tantai-kanban/kanban-api/internal/cards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

func (r implRepository) buildModel(ctx context.Context, opts repository.CreateOptions) dbmodels.Card {
	m := dbmodels.Card{
		ListID:      opts.ListID,
		Name:        opts.Name,
		Description: null.StringFrom(opts.Description),
		Position:    types.Decimal{Big: decimal.New(int64(opts.Position), 0)},
		Priority:    dbmodels.CardPriority(opts.Priority),
		DueDate:     null.TimeFromPtr(opts.DueDate),
		CreatedBy:   null.StringFrom(opts.CreatedBy),
		CreatedAt:   r.clock(),
		UpdatedAt:   r.clock(),
	}

	// Convert labels to JSON
	if len(opts.Labels) > 0 {
		labelsJSON, _ := json.Marshal(opts.Labels)
		m.Labels = null.JSONFrom(labelsJSON)
	}

	// Handle new fields
	if opts.AssignedTo != nil && *opts.AssignedTo != "" {
		m.AssignedTo = null.StringFrom(*opts.AssignedTo)
	}

	if opts.EstimatedHours != nil {
		// TODO: Implement proper decimal conversion
		// For now, skip estimated hours to avoid decimal issues
	}

	if opts.StartDate != nil {
		m.StartDate = null.TimeFromPtr(opts.StartDate)
	}

	if len(opts.Tags) > 0 {
		m.Tags = opts.Tags
	}

	if len(opts.Checklist) > 0 {
		checklistJSON, _ := json.Marshal(opts.Checklist)
		m.Checklist = null.JSONFrom(checklistJSON)
	}

	return m
}

func (r implRepository) buildUpdateModel(ctx context.Context, opts repository.UpdateOptions) (dbmodels.Card, []string, map[string]interface{}, error) {
	card := dbmodels.Card{}
	cols := make([]string, 0)
	updates := make(map[string]interface{})

	if opts.Name != nil {
		card.Name = *opts.Name
		cols = append(cols, dbmodels.CardColumns.Name)
		updates["name"] = *opts.Name
	}
	if opts.Description != nil {
		card.Description = null.StringFrom(*opts.Description)
		cols = append(cols, dbmodels.CardColumns.Description)
		updates["description"] = *opts.Description
	}
	if opts.Priority != nil {
		card.Priority = dbmodels.CardPriority(*opts.Priority)
		cols = append(cols, dbmodels.CardColumns.Priority)
		updates["priority"] = *opts.Priority
	}
	if opts.Labels != nil {
		labelsJSON, _ := json.Marshal(*opts.Labels)
		card.Labels = null.JSONFrom(labelsJSON)
		cols = append(cols, dbmodels.CardColumns.Labels)
		updates["labels"] = *opts.Labels
	}
	if opts.DueDate != nil {
		card.DueDate = null.TimeFromPtr(opts.DueDate)
		cols = append(cols, dbmodels.CardColumns.DueDate)
		updates["due_date"] = opts.DueDate
	}

	// Handle new fields
	if opts.AssignedTo != nil {
		if *opts.AssignedTo != "" {
			card.AssignedTo = null.StringFrom(*opts.AssignedTo)
		} else {
			card.AssignedTo = null.StringFromPtr(nil) // Set to NULL when empty string
		}
		cols = append(cols, dbmodels.CardColumns.AssignedTo)
		updates["assigned_to"] = *opts.AssignedTo
	}

	if opts.EstimatedHours != nil {
		// TODO: Implement proper decimal conversion
		// For now, skip estimated hours to avoid decimal issues
	}

	if opts.ActualHours != nil {
		// TODO: Implement proper decimal conversion
		// For now, skip actual hours to avoid decimal issues
	}

	if opts.StartDate != nil {
		card.StartDate = null.TimeFromPtr(opts.StartDate)
		cols = append(cols, dbmodels.CardColumns.StartDate)
		updates["start_date"] = opts.StartDate
	}

	if opts.CompletionDate != nil {
		card.CompletionDate = null.TimeFromPtr(opts.CompletionDate)
		cols = append(cols, dbmodels.CardColumns.CompletionDate)
		updates["completion_date"] = opts.CompletionDate
	}

	if opts.Tags != nil {
		card.Tags = *opts.Tags
		cols = append(cols, dbmodels.CardColumns.Tags)
		updates["tags"] = *opts.Tags
	}

	if opts.Checklist != nil {
		checklistJSON, _ := json.Marshal(*opts.Checklist)
		card.Checklist = null.JSONFrom(checklistJSON)
		cols = append(cols, dbmodels.CardColumns.Checklist)
		updates["checklist"] = *opts.Checklist
	}

	if err := postgres.IsUUID(opts.ID); err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.buildUpdateModel.IsUUID: %v", err)
		return dbmodels.Card{}, nil, nil, err
	}
	card.ID = opts.ID

	card.UpdatedAt = r.clock()
	cols = append(cols, dbmodels.CardColumns.UpdatedAt)

	return card, cols, updates, nil
}

func (r implRepository) buildMoveModel(ctx context.Context, opts repository.MoveOptions) (dbmodels.Card, []string, error) {
	card := dbmodels.Card{
		ListID:   opts.ListID,
		Position: types.Decimal{Big: decimal.New(int64(opts.Position), 0)},
	}
	cols := make([]string, 0)
	cols = append(cols, dbmodels.CardColumns.ListID)
	cols = append(cols, dbmodels.CardColumns.Position)

	if err := postgres.IsUUID(opts.ID); err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.buildMoveModel.IsUUID: %v", err)
		return dbmodels.Card{}, nil, err
	}
	card.ID = opts.ID

	card.UpdatedAt = r.clock()
	cols = append(cols, dbmodels.CardColumns.UpdatedAt)

	return card, cols, nil
}

func (r implRepository) buildActivityModel(ctx context.Context, cardID string, actionType string, oldData, newData map[string]interface{}) dbmodels.CardActivity {
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
