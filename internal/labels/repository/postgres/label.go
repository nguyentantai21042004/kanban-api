package postgres

import (
	"context"
	"database/sql"
	"sync"

	"github.com/aarondl/sqlboiler/v4/boil"
	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/internal/labels/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

func (r implRepository) Get(ctx context.Context, sc models.Scope, opts repository.GetOptions) ([]models.Label, paginator.Paginator, error) {
	qr, err := r.buildGetQuery(ctx, opts.Filter)
	if err != nil {
		r.l.Errorf(ctx, "internal.labels.repository.postgres.Get.buildGetQuery: %v", err)
		return nil, paginator.Paginator{}, err
	}

	var (
		total int64
		ls    dbmodels.LabelSlice
	)

	errChan := make(chan error, 2)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		var countErr error
		total, countErr = dbmodels.Labels(qr...).Count(ctx, r.database)
		if countErr != nil {
			r.l.Errorf(ctx, "internal.labels.repository.postgres.Get.Count: %v", countErr)
			errChan <- countErr
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ls, err = dbmodels.Labels(qr...).All(ctx, r.database)
		if err != nil {
			r.l.Errorf(ctx, "internal.labels.repository.postgres.Get.All: %v", err)
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			r.l.Errorf(ctx, "internal.labels.repository.postgres.Get.errChan: %v", err)
			return nil, paginator.Paginator{}, err
		}
	}
	dbLabels := util.DerefSlice(ls)
	labels := make([]models.Label, len(dbLabels))
	for i, label := range dbLabels {
		labels[i] = models.NewLabel(label)
	}

	return labels, paginator.Paginator{
		Total:       total,
		Count:       int64(len(labels)),
		PerPage:     opts.PagQuery.Limit,
		CurrentPage: opts.PagQuery.Page,
	}, nil
}

func (r implRepository) Create(ctx context.Context, sc models.Scope, opts repository.CreateOptions) (models.Label, error) {
	m := r.buildModel(ctx, opts)

	err := m.Insert(ctx, r.database, boil.Infer())
	if err != nil {
		r.l.Errorf(ctx, "internal.labels.repository.postgres.Create.Insert: %v", err)
		return models.Label{}, err
	}

	return models.NewLabel(m), nil
}

func (r implRepository) Update(ctx context.Context, sc models.Scope, opts repository.UpdateOptions) (models.Label, error) {
	l, col, err := r.buildUpdateModel(ctx, opts)
	if err != nil {
		r.l.Errorf(ctx, "internal.labels.repository.postgres.Update.buildUpdateModel: %v", err)
		return models.Label{}, err
	}

	_, err = l.Update(ctx, r.database, boil.Whitelist(col...))
	if err != nil {
		r.l.Errorf(ctx, "internal.labels.repository.postgres.Update.Update: %v", err)
		return models.Label{}, err
	}

	return models.NewLabel(l), nil
}

func (r implRepository) Detail(ctx context.Context, sc models.Scope, ID string) (models.Label, error) {
	l, err := dbmodels.Labels(dbmodels.LabelWhere.ID.EQ(ID)).One(ctx, r.database)
	if err != nil {
		if err == sql.ErrNoRows {
			r.l.Warnf(ctx, "internal.labels.repository.postgres.Detail.One.NotFound: %v", err)
			return models.Label{}, repository.ErrNotFound
		}
		r.l.Errorf(ctx, "internal.labels.repository.postgres.Detail.One: %v", err)
		return models.Label{}, err
	}

	return models.NewLabel(*l), nil
}

func (r implRepository) Delete(ctx context.Context, sc models.Scope, IDs []string) error {
	_, err := dbmodels.Labels(dbmodels.LabelWhere.ID.IN(IDs)).DeleteAll(ctx, r.database, false)
	if err != nil {
		r.l.Errorf(ctx, "internal.labels.repository.postgres.Delete.DeleteAll: %v", err)
		return err
	}

	return nil
}
