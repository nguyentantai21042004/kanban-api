package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync"

	"github.com/aarondl/sqlboiler/v4/boil"
	"gitlab.com/tantai-kanban/kanban-api/internal/cards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

func (r implRepository) Get(ctx context.Context, sc models.Scope, opts repository.GetOptions) ([]models.Card, paginator.Paginator, error) {
	qr, err := r.buildGetQuery(ctx, opts.Filter)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Get.buildGetQuery: %v", err)
		return nil, paginator.Paginator{}, err
	}

	var (
		total int64
		cs    dbmodels.CardSlice
	)

	errChan := make(chan error, 2)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		var countErr error
		total, countErr = dbmodels.Cards(qr...).Count(ctx, r.database)
		if countErr != nil {
			r.l.Errorf(ctx, "internal.cards.repository.postgres.Get.Count: %v", countErr)
			errChan <- countErr
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cs, err = dbmodels.Cards(qr...).All(ctx, r.database)
		if err != nil {
			r.l.Errorf(ctx, "internal.cards.repository.postgres.Get.All: %v", err)
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			r.l.Errorf(ctx, "internal.cards.repository.postgres.Get.errChan: %v", err)
			return nil, paginator.Paginator{}, err
		}
	}
	dbCards := util.DerefSlice(cs)
	cards := make([]models.Card, len(dbCards))
	for i, card := range dbCards {
		cards[i] = models.NewCard(card)
	}

	return cards, paginator.Paginator{
		Total:       total,
		Count:       int64(len(cards)),
		PerPage:     opts.PagQuery.Limit,
		CurrentPage: opts.PagQuery.Page,
	}, nil
}

func (r implRepository) Create(ctx context.Context, sc models.Scope, opts repository.CreateOptions) (models.Card, error) {
	// Start transaction
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Create.BeginTx: %v", err)
		return models.Card{}, err
	}
	defer tx.Rollback()

	m := r.buildModel(ctx, opts)

	err = m.Insert(ctx, tx, boil.Infer())
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Create.Insert: %v", err)
		return models.Card{}, err
	}

	// Create activity record
	activity := r.buildActivityModel(ctx, m.ID, string(models.CardActionTypeCreated), nil, map[string]interface{}{
		"Name":        m.Name,
		"description": m.Description,
		"priority":    m.Priority,
		"labels":      m.Labels,
	})

	err = activity.Insert(ctx, tx, boil.Infer())
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Create.InsertActivity: %v", err)
		return models.Card{}, err
	}

	if err := tx.Commit(); err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Create.Commit: %v", err)
		return models.Card{}, err
	}

	return models.NewCard(m), nil
}

func (r implRepository) Update(ctx context.Context, sc models.Scope, opts repository.UpdateOptions) (models.Card, error) {
	// Start transaction
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Update.BeginTx: %v", err)
		return models.Card{}, err
	}
	defer tx.Rollback()

	c, col, updates, err := r.buildUpdateModel(ctx, opts)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Update.buildUpdateModel: %v", err)
		return models.Card{}, err
	}

	_, err = c.Update(ctx, tx, boil.Whitelist(col...))
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Update.Update: %v", err)
		return models.Card{}, err
	}

	// Create activity record if there were updates
	if len(updates) > 0 {
		oldData := map[string]interface{}{
			"Name":        opts.OldModel.Name,
			"description": opts.OldModel.Description,
			"priority":    opts.OldModel.Priority,
			"labels":      opts.OldModel.Labels,
		}
		activity := r.buildActivityModel(ctx, c.ID, string(models.CardActionTypeUpdated), oldData, updates)
		err = activity.Insert(ctx, tx, boil.Infer())
		if err != nil {
			r.l.Errorf(ctx, "internal.cards.repository.postgres.Update.InsertActivity: %v", err)
			return models.Card{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Update.Commit: %v", err)
		return models.Card{}, err
	}

	return models.NewCard(c), nil
}

// DEPRECATED: Use EnhancedMove with position manager instead
// recalculatePositions is deprecated in favor of position manager
func (r implRepository) recalculatePositions(ctx context.Context, listID string, excludeCardID string) error {
	// This method is deprecated. Use RebalanceListPositions instead
	return r.RebalanceListPositions(ctx, listID)
}

// DEPRECATED: Use position manager instead
// calculateNewPosition is deprecated in favor of position manager
func (r implRepository) calculateNewPosition(ctx context.Context, listID string, targetPosition float64) (float64, error) {
	// This method is deprecated. Use position manager's GeneratePosition instead
	return targetPosition, nil
}

// DEPRECATED: Use EnhancedMove instead
// Move is deprecated in favor of EnhancedMove with advanced position management
func (r implRepository) Move(ctx context.Context, sc models.Scope, opts repository.MoveOptions) (models.Card, error) {
	// Start transaction
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Move.BeginTx: %v", err)
		return models.Card{}, err
	}
	defer tx.Rollback()

	// Tính toán position mới nếu cần
	newPosition := opts.Position
	if opts.Position <= 0 {
		calculatedPos, err := r.calculateNewPosition(ctx, opts.ListID, 0)
		if err != nil {
			r.l.Errorf(ctx, "internal.cards.repository.postgres.Move.calculateNewPosition: %v", err)
			return models.Card{}, err
		}
		newPosition = calculatedPos
	}

	// Nếu move trong cùng list, tính toán lại position của các card khác
	if opts.OldModel.ListID == opts.ListID {
		err = r.recalculatePositions(ctx, opts.ListID, opts.ID)
		if err != nil {
			r.l.Errorf(ctx, "internal.cards.repository.postgres.Move.recalculatePositions: %v", err)
			return models.Card{}, err
		}
	}

	c, col, err := r.buildMoveModel(ctx, repository.MoveOptions{
		ID:       opts.ID,
		ListID:   opts.ListID,
		Position: newPosition,
		OldModel: opts.OldModel,
	})
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Move.buildMoveModel: %v", err)
		return models.Card{}, err
	}

	_, err = c.Update(ctx, tx, boil.Whitelist(col...))
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Move.Update: %v", err)
		return models.Card{}, err
	}

	// Create activity record
	activity := r.buildActivityModel(ctx, c.ID, string(models.CardActionTypeMoved), map[string]interface{}{
		"list_id":  opts.OldModel.ListID,
		"position": opts.OldModel.Position,
	}, map[string]interface{}{
		"list_id":  opts.ListID,
		"position": newPosition,
	})

	err = activity.Insert(ctx, tx, boil.Infer())
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Move.InsertActivity: %v", err)
		return models.Card{}, err
	}

	if err := tx.Commit(); err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Move.Commit: %v", err)
		return models.Card{}, err
	}

	return models.NewCard(c), nil
}

func (r implRepository) Detail(ctx context.Context, sc models.Scope, ID string) (models.Card, error) {
	c, err := dbmodels.Cards(dbmodels.CardWhere.ID.EQ(ID)).One(ctx, r.database)
	if err != nil {
		if err == sql.ErrNoRows {
			r.l.Warnf(ctx, "internal.cards.repository.postgres.Detail.One.NotFound: %v", err)
			return models.Card{}, repository.ErrNotFound
		}
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Detail.One: %v", err)
		return models.Card{}, err
	}

	return models.NewCard(*c), nil
}

func (r implRepository) Delete(ctx context.Context, sc models.Scope, IDs []string) error {
	_, err := dbmodels.Cards(dbmodels.CardWhere.ID.IN(IDs)).DeleteAll(ctx, r.database, true)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Delete.DeleteAll: %v", err)
		return err
	}

	return nil
}

func (r implRepository) GetMaxPosition(ctx context.Context, sc models.Scope, listID string) (float64, error) {
	var maxPosition float64
	cards, err := dbmodels.Cards(dbmodels.CardWhere.ListID.EQ(listID)).All(ctx, r.database)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.GetMaxPosition.All: %v", err)
		return 0, err
	}

	if len(cards) == 0 {
		return 0, nil
	}

	// Find max position
	for _, card := range cards {
		if card.Position.Big != nil {
			pos, _ := card.Position.Big.Float64()
			if pos > maxPosition {
				maxPosition = pos
			}
		}
	}

	return maxPosition, nil
}

func (r implRepository) GetActivities(ctx context.Context, sc models.Scope, opts repository.GetActivitiesOptions) ([]models.CardActivity, error) {
	activities, err := dbmodels.CardActivities(
		dbmodels.CardActivityWhere.CardID.EQ(opts.CardID),
	).All(ctx, r.database)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.GetActivities.All: %v", err)
		return nil, err
	}

	dbActivities := util.DerefSlice(activities)
	result := make([]models.CardActivity, len(dbActivities))
	for i, activity := range dbActivities {
		result[i] = models.NewCardActivity(activity)
	}

	return result, nil
}

// GetBoardIDFromListID retrieves the board ID for a given list ID
func (r implRepository) GetBoardIDFromListID(ctx context.Context, listID string) (string, error) {
	list, err := dbmodels.FindList(ctx, r.database, listID)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.GetBoardIDFromListID.FindList: %v", err)
		return "", err
	}
	return list.BoardID, nil
}

// New methods for enhanced functionality
func (r implRepository) Assign(ctx context.Context, sc models.Scope, opts repository.AssignOptions) (models.Card, error) {
	card, err := dbmodels.FindCard(ctx, r.database, opts.CardID)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Assign.FindCard: %v", err)
		return models.Card{}, err
	}

	card.AssignedTo.String = opts.AssignedTo
	card.AssignedTo.Valid = true
	card.UpdatedAt = r.clock()

	_, err = card.Update(ctx, r.database, boil.Whitelist(dbmodels.CardColumns.AssignedTo, dbmodels.CardColumns.UpdatedAt))
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Assign.Update: %v", err)
		return models.Card{}, err
	}

	return models.NewCard(*card), nil
}

func (r implRepository) Unassign(ctx context.Context, sc models.Scope, opts repository.UnassignOptions) (models.Card, error) {
	card, err := dbmodels.FindCard(ctx, r.database, opts.CardID)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Unassign.FindCard: %v", err)
		return models.Card{}, err
	}

	card.AssignedTo.Valid = false
	card.UpdatedAt = r.clock()

	_, err = card.Update(ctx, r.database, boil.Whitelist(dbmodels.CardColumns.AssignedTo, dbmodels.CardColumns.UpdatedAt))
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Unassign.Update: %v", err)
		return models.Card{}, err
	}

	return models.NewCard(*card), nil
}

func (r implRepository) AddAttachment(ctx context.Context, sc models.Scope, opts repository.AddAttachmentOptions) (models.Card, error) {
	card, err := dbmodels.FindCard(ctx, r.database, opts.CardID)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.AddAttachment.FindCard: %v", err)
		return models.Card{}, err
	}

	attachments := []string{}
	if card.Attachments.Valid {
		_ = json.Unmarshal(card.Attachments.JSON, &attachments)
	}

	// Check if attachment already exists
	for _, attachment := range attachments {
		if attachment == opts.AttachmentID {
			return models.NewCard(*card), nil // Already exists
		}
	}

	attachments = append(attachments, opts.AttachmentID)
	attachmentsJSON, _ := json.Marshal(attachments)
	card.Attachments.JSON = attachmentsJSON
	card.Attachments.Valid = true
	card.UpdatedAt = r.clock()

	_, err = card.Update(ctx, r.database, boil.Whitelist(dbmodels.CardColumns.Attachments, dbmodels.CardColumns.UpdatedAt))
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.AddAttachment.Update: %v", err)
		return models.Card{}, err
	}

	return models.NewCard(*card), nil
}

func (r implRepository) RemoveAttachment(ctx context.Context, sc models.Scope, opts repository.RemoveAttachmentOptions) (models.Card, error) {
	card, err := dbmodels.FindCard(ctx, r.database, opts.CardID)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.RemoveAttachment.FindCard: %v", err)
		return models.Card{}, err
	}

	attachments := []string{}
	if card.Attachments.Valid {
		_ = json.Unmarshal(card.Attachments.JSON, &attachments)
	}

	// Remove attachment
	newAttachments := []string{}
	for _, attachment := range attachments {
		if attachment != opts.AttachmentID {
			newAttachments = append(newAttachments, attachment)
		}
	}

	attachmentsJSON, _ := json.Marshal(newAttachments)
	card.Attachments.JSON = attachmentsJSON
	card.Attachments.Valid = true
	card.UpdatedAt = r.clock()

	_, err = card.Update(ctx, r.database, boil.Whitelist(dbmodels.CardColumns.Attachments, dbmodels.CardColumns.UpdatedAt))
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.RemoveAttachment.Update: %v", err)
		return models.Card{}, err
	}

	return models.NewCard(*card), nil
}

func (r implRepository) UpdateTimeTracking(ctx context.Context, sc models.Scope, opts repository.UpdateTimeTrackingOptions) (models.Card, error) {
	card, err := dbmodels.FindCard(ctx, r.database, opts.CardID)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.UpdateTimeTracking.FindCard: %v", err)
		return models.Card{}, err
	}

	// TODO: Implement time tracking update
	// For now, just update the timestamp
	card.UpdatedAt = r.clock()

	_, err = card.Update(ctx, r.database, boil.Whitelist(dbmodels.CardColumns.UpdatedAt))
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.UpdateTimeTracking.Update: %v", err)
		return models.Card{}, err
	}

	return models.NewCard(*card), nil
}

func (r implRepository) UpdateChecklist(ctx context.Context, sc models.Scope, opts repository.UpdateChecklistOptions) (models.Card, error) {
	card, err := dbmodels.FindCard(ctx, r.database, opts.CardID)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.UpdateChecklist.FindCard: %v", err)
		return models.Card{}, err
	}

	checklistJSON, _ := json.Marshal(opts.Checklist)
	card.Checklist.JSON = checklistJSON
	card.Checklist.Valid = true
	card.UpdatedAt = r.clock()

	_, err = card.Update(ctx, r.database, boil.Whitelist(dbmodels.CardColumns.Checklist, dbmodels.CardColumns.UpdatedAt))
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.UpdateChecklist.Update: %v", err)
		return models.Card{}, err
	}

	return models.NewCard(*card), nil
}

func (r implRepository) AddTag(ctx context.Context, sc models.Scope, opts repository.AddTagOptions) (models.Card, error) {
	card, err := dbmodels.FindCard(ctx, r.database, opts.CardID)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.AddTag.FindCard: %v", err)
		return models.Card{}, err
	}

	tags := card.Tags
	// Check if tag already exists
	for _, tag := range tags {
		if tag == opts.Tag {
			return models.NewCard(*card), nil // Already exists
		}
	}

	tags = append(tags, opts.Tag)
	card.Tags = tags
	card.UpdatedAt = r.clock()

	_, err = card.Update(ctx, r.database, boil.Whitelist(dbmodels.CardColumns.Tags, dbmodels.CardColumns.UpdatedAt))
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.AddTag.Update: %v", err)
		return models.Card{}, err
	}

	return models.NewCard(*card), nil
}

func (r implRepository) RemoveTag(ctx context.Context, sc models.Scope, opts repository.RemoveTagOptions) (models.Card, error) {
	card, err := dbmodels.FindCard(ctx, r.database, opts.CardID)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.RemoveTag.FindCard: %v", err)
		return models.Card{}, err
	}

	tags := card.Tags
	newTags := []string{}
	for _, tag := range tags {
		if tag != opts.Tag {
			newTags = append(newTags, tag)
		}
	}

	card.Tags = newTags
	card.UpdatedAt = r.clock()

	_, err = card.Update(ctx, r.database, boil.Whitelist(dbmodels.CardColumns.Tags, dbmodels.CardColumns.UpdatedAt))
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.RemoveTag.Update: %v", err)
		return models.Card{}, err
	}

	return models.NewCard(*card), nil
}

func (r implRepository) SetStartDate(ctx context.Context, sc models.Scope, opts repository.SetStartDateOptions) (models.Card, error) {
	card, err := dbmodels.FindCard(ctx, r.database, opts.CardID)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.SetStartDate.FindCard: %v", err)
		return models.Card{}, err
	}

	if opts.StartDate != nil {
		card.StartDate.Time = *opts.StartDate
		card.StartDate.Valid = true
	} else {
		card.StartDate.Valid = false
	}

	card.UpdatedAt = r.clock()

	_, err = card.Update(ctx, r.database, boil.Whitelist(dbmodels.CardColumns.StartDate, dbmodels.CardColumns.UpdatedAt))
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.SetStartDate.Update: %v", err)
		return models.Card{}, err
	}

	return models.NewCard(*card), nil
}

func (r implRepository) SetCompletionDate(ctx context.Context, sc models.Scope, opts repository.SetCompletionDateOptions) (models.Card, error) {
	card, err := dbmodels.FindCard(ctx, r.database, opts.CardID)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.SetCompletionDate.FindCard: %v", err)
		return models.Card{}, err
	}

	if opts.CompletionDate != nil {
		card.CompletionDate.Time = *opts.CompletionDate
		card.CompletionDate.Valid = true
	} else {
		card.CompletionDate.Valid = false
	}

	card.UpdatedAt = r.clock()

	_, err = card.Update(ctx, r.database, boil.Whitelist(dbmodels.CardColumns.CompletionDate, dbmodels.CardColumns.UpdatedAt))
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.SetCompletionDate.Update: %v", err)
		return models.Card{}, err
	}

	return models.NewCard(*card), nil
}
