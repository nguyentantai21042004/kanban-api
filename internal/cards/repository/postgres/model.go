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

func (r implRepository) buildModel(ctx context.Context, opts repository.CreateOptions) (dbmodels.Card, error) {
	if opts.BoardID == "" {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Create.Validation: BoardID is required")
		return dbmodels.Card{}, repository.ErrFieldRequired
	}
	if opts.ListID == "" {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Create.Validation: ListID is required")
		return dbmodels.Card{}, repository.ErrFieldRequired
	}
	if opts.Name == "" {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Create.Validation: Name is required")
		return dbmodels.Card{}, repository.ErrFieldRequired
	}
	if opts.CreatedBy == "" {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Create.Validation: CreatedBy is required")
		return dbmodels.Card{}, repository.ErrFieldRequired
	}

	m := dbmodels.Card{
		BoardID:     opts.BoardID,
		ListID:      opts.ListID,
		Name:        opts.Name,
		Alias:       null.StringFrom(opts.Alias),
		Description: null.StringFrom(opts.Description),
		Position:    opts.Position,
		Priority:    dbmodels.CardPriority(opts.Priority),
		DueDate:     null.TimeFromPtr(opts.DueDate),
		CreatedBy:   null.StringFrom(opts.CreatedBy),
		CreatedAt:   r.clock(),
		UpdatedAt:   r.clock(),
	}

	if len(opts.Labels) > 0 {
		labelsJSON, _ := json.Marshal(opts.Labels)
		m.Labels = null.JSONFrom(labelsJSON)
	}

	if opts.AssignedTo != nil && *opts.AssignedTo != "" {
		m.AssignedTo = null.StringFrom(*opts.AssignedTo)
	}

	if opts.EstimatedHours != nil {
		m.EstimatedHours = types.NullDecimal{Big: decimal.New(int64(*opts.EstimatedHours), 0)}
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

	return m, nil
}

func (r implRepository) buildUpdateModel(ctx context.Context, opts repository.UpdateOptions) (dbmodels.Card, []string, map[string]interface{}, error) {
	card := dbmodels.Card{}
	cols := make([]string, 0)
	updates := make(map[string]interface{})

	if opts.Name != "" {
		card.Name = opts.Name
		cols = append(cols, dbmodels.CardColumns.Name)
		updates["name"] = opts.Name
	}
	if opts.Alias != "" {
		card.Alias = null.StringFrom(opts.Alias)
		cols = append(cols, dbmodels.CardColumns.Alias)
		updates["alias"] = opts.Alias
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
		card.EstimatedHours = types.NullDecimal{Big: decimal.New(int64(*opts.EstimatedHours), 0)}
		cols = append(cols, dbmodels.CardColumns.EstimatedHours)
		updates["estimated_hours"] = *opts.EstimatedHours
	}

	if opts.ActualHours != nil {
		card.ActualHours = types.NullDecimal{Big: decimal.New(int64(*opts.ActualHours), 0)}
		cols = append(cols, dbmodels.CardColumns.ActualHours)
		updates["actual_hours"] = *opts.ActualHours
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
		Position: opts.NewPosition,
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
