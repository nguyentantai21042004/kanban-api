package postgres

import (
	"context"
	"database/sql"
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
		"title":       m.Title,
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
			"title":       opts.OldModel.Title,
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

func (r implRepository) Move(ctx context.Context, sc models.Scope, opts repository.MoveOptions) (models.Card, error) {
	// Start transaction
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		r.l.Errorf(ctx, "internal.cards.repository.postgres.Move.BeginTx: %v", err)
		return models.Card{}, err
	}
	defer tx.Rollback()

	c, col, err := r.buildMoveModel(ctx, opts)
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
		"position": opts.Position,
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
