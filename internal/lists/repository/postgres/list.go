package postgres

import (
	"context"
	"database/sql"
	"sync"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/nguyentantai21042004/kanban-api/internal/dbmodels"
	"github.com/nguyentantai21042004/kanban-api/internal/lists/repository"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/pkg/paginator"
	"github.com/nguyentantai21042004/kanban-api/pkg/util"
)

func (r implRepository) GetPosition(ctx context.Context, sc models.Scope, opts repository.GetPositionOptions) (string, error) {
	qr, err := r.buildGetPositionQuery(opts)
	if err != nil {
		r.l.Errorf(ctx, "internal.lists.repository.postgres.GetPosition.buildGetPositionQuery: %v", err)
		return "", err
	}

	pos, err := dbmodels.Lists(qr...).One(ctx, r.database)
	if err != nil {
		if err == sql.ErrNoRows {
			r.l.Warnf(ctx, "internal.lists.repository.postgres.GetPosition.One.NotFound: %v", err)
			return "", repository.ErrNotFound
		}
		r.l.Errorf(ctx, "internal.lists.repository.postgres.GetPosition.One: %v", err)
		return "", err
	}

	return pos.Position, nil
}

func (r implRepository) List(ctx context.Context, sc models.Scope, opts repository.ListOptions) ([]models.List, error) {
	qr, err := r.buildGetQuery(ctx, opts.Filter)
	if err != nil {
		r.l.Errorf(ctx, "internal.lists.repository.postgres.List.buildListQuery: %v", err)
		return nil, err
	}

	ls, err := dbmodels.Lists(qr...).All(ctx, r.database)
	if err != nil {
		r.l.Errorf(ctx, "internal.lists.repository.postgres.List.All: %v", err)
		return nil, err
	}

	dbLists := util.DerefSlice(ls)
	lists := make([]models.List, len(dbLists))
	for i, list := range dbLists {
		lists[i] = models.NewList(list)
	}

	return lists, nil
}

func (r implRepository) Get(ctx context.Context, sc models.Scope, opts repository.GetOptions) ([]models.List, paginator.Paginator, error) {
	qr, err := r.buildGetQuery(ctx, opts.Filter)
	if err != nil {
		r.l.Errorf(ctx, "internal.lists.repository.postgres.Get.buildGetQuery: %v", err)
		return nil, paginator.Paginator{}, err
	}

	var (
		total int64
		ls    dbmodels.ListSlice
	)

	errChan := make(chan error, 2)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		var countErr error
		total, countErr = dbmodels.Lists(qr...).Count(ctx, r.database)
		if countErr != nil {
			r.l.Errorf(ctx, "internal.lists.repository.postgres.Get.Count: %v", countErr)
			errChan <- countErr
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ls, err = dbmodels.Lists(qr...).All(ctx, r.database)
		if err != nil {
			r.l.Errorf(ctx, "internal.lists.repository.postgres.Get.All: %v", err)
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			r.l.Errorf(ctx, "internal.lists.repository.postgres.Get.errChan: %v", err)
			return nil, paginator.Paginator{}, err
		}
	}
	dbLists := util.DerefSlice(ls)
	lists := make([]models.List, len(dbLists))
	for i, list := range dbLists {
		lists[i] = models.NewList(list)
	}

	return lists, paginator.Paginator{
		Total:       total,
		Count:       int64(len(lists)),
		PerPage:     opts.PagQuery.Limit,
		CurrentPage: opts.PagQuery.Page,
	}, nil
}

func (r implRepository) Create(ctx context.Context, sc models.Scope, opts repository.CreateOptions) (models.List, error) {
	m := r.buildModel(ctx, sc, opts)

	err := m.Insert(ctx, r.database, boil.Infer())
	if err != nil {
		r.l.Errorf(ctx, "internal.lists.repository.postgres.Create.Insert: %v", err)
		return models.List{}, err
	}

	return models.NewList(m), nil
}

func (r implRepository) Update(ctx context.Context, sc models.Scope, opts repository.UpdateOptions) (models.List, error) {
	l, col, err := r.buildUpdateModel(ctx, opts)
	if err != nil {
		r.l.Errorf(ctx, "internal.lists.repository.postgres.Update.buildUpdateModel: %v", err)
		return models.List{}, err
	}

	_, err = l.Update(ctx, r.database, boil.Whitelist(col...))
	if err != nil {
		r.l.Errorf(ctx, "internal.lists.repository.postgres.Update.Update: %v", err)
		return models.List{}, err
	}

	return models.NewList(l), nil
}

func (r implRepository) Detail(ctx context.Context, sc models.Scope, ID string) (models.List, error) {
	l, err := dbmodels.Lists(dbmodels.ListWhere.ID.EQ(ID)).One(ctx, r.database)
	if err != nil {
		if err == sql.ErrNoRows {
			r.l.Warnf(ctx, "internal.lists.repository.postgres.Detail.One.NotFound: %v", err)
			return models.List{}, repository.ErrNotFound
		}
		r.l.Errorf(ctx, "internal.lists.repository.postgres.Detail.One: %v", err)
		return models.List{}, err
	}

	return models.NewList(*l), nil
}

func (r implRepository) Delete(ctx context.Context, sc models.Scope, IDs []string) error {
	_, err := dbmodels.Lists(dbmodels.ListWhere.ID.IN(IDs)).DeleteAll(ctx, r.database, true)
	if err != nil {
		r.l.Errorf(ctx, "internal.lists.repository.postgres.Delete.DeleteAll: %v", err)
		return err
	}

	return nil
}

func (r implRepository) Move(ctx context.Context, sc models.Scope, opts repository.MoveOptions) (models.List, error) {
	l, err := dbmodels.FindList(ctx, r.database, opts.ID)
	if err != nil {
		r.l.Errorf(ctx, "internal.lists.repository.postgres.Move.FindList: %v", err)
		return models.List{}, err
	}

	l.Position = opts.NewPosition
	l.UpdatedAt = r.clock()

	_, err = l.Update(ctx, r.database, boil.Whitelist(dbmodels.ListColumns.Position, dbmodels.ListColumns.UpdatedAt))
	if err != nil {
		r.l.Errorf(ctx, "internal.lists.repository.postgres.Move.Update: %v", err)
		return models.List{}, err
	}

	return models.NewList(*l), nil
}
