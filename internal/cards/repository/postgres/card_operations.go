package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/nguyentantai21042004/kanban-api/internal/cards/repository"
	"github.com/nguyentantai21042004/kanban-api/internal/dbmodels"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/pkg/paginator"
	"github.com/nguyentantai21042004/kanban-api/pkg/util"
)

func (r implRepository) GetPosition(ctx context.Context, sc models.Scope, opts repository.GetPositionOptions) (string, error) {
	qr, err := r.buildGetPositionQuery(opts)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.GetPosition.buildGetPositionQuery: %v", err)
		return "", err
	}

	card, err := dbmodels.Cards(qr...).One(ctx, r.database)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", repository.ErrNotFound
		}
		r.l.Errorf(ctx, "internal.cards.repository.postgres.GetPosition.One: %v", err)
		return "", err
	}

	return card.Position, nil
}

func (r implRepository) GetActivities(ctx context.Context, sc models.Scope, opts repository.GetActivitiesOptions) ([]models.CardActivity, paginator.Paginator, error) {
	qr, err := r.buildGetActivitiesQuery(ctx, opts.Filter)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.GetActivities.buildGetActivitiesQuery: %v", err)
		return nil, paginator.Paginator{}, err
	}

	var (
		total int64
		cs    dbmodels.CardActivitySlice
	)

	errChan := make(chan error, 2)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		var countErr error
		total, countErr = dbmodels.CardActivities(qr...).Count(ctx, r.database)
		if countErr != nil {
			r.l.Errorf(ctx, "internal.cards.repository.postgres.GetActivities.Count: %v", countErr)
			errChan <- countErr
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cs, err = dbmodels.CardActivities(qr...).All(ctx, r.database)
		if err != nil {
			r.l.Errorf(ctx, "internal.cards.repository.postgres.GetActivities.All: %v", err)
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			r.l.Errorf(ctx, "internal.cards.repository.postgres.GetActivities.errChan: %v", err)
			return nil, paginator.Paginator{}, err
		}
	}
	dbActivities := util.DerefSlice(cs)
	activities := make([]models.CardActivity, len(dbActivities))
	for i, activity := range dbActivities {
		activities[i] = models.NewCardActivity(activity)
	}

	return activities, paginator.Paginator{
		Total:       total,
		Count:       int64(len(activities)),
		PerPage:     opts.PagQuery.Limit,
		CurrentPage: opts.PagQuery.Page,
	}, nil
}

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
