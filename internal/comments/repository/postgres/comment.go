package postgres

import (
	"context"
	"database/sql"
	"sync"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/nguyentantai21042004/kanban-api/internal/comments/repository"
	"github.com/nguyentantai21042004/kanban-api/internal/dbmodels"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/pkg/paginator"
	"github.com/nguyentantai21042004/kanban-api/pkg/util"
)

func (r implRepository) Get(ctx context.Context, sc models.Scope, opts repository.GetOptions) ([]models.Comment, paginator.Paginator, error) {
	qr, err := r.buildGetQuery(ctx, opts.Filter)
	if err != nil {
		r.l.Errorf(ctx, "internal.comments.repository.postgres.Get.buildGetQuery: %v", err)
		return nil, paginator.Paginator{}, err
	}

	var (
		total int64
		cs    dbmodels.CommentSlice
	)

	errChan := make(chan error, 2)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		var countErr error
		total, countErr = dbmodels.Comments(qr...).Count(ctx, r.database)
		if countErr != nil {
			r.l.Errorf(ctx, "internal.comments.repository.postgres.Get.Count: %v", countErr)
			errChan <- countErr
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cs, err = dbmodels.Comments(qr...).All(ctx, r.database)
		if err != nil {
			r.l.Errorf(ctx, "internal.comments.repository.postgres.Get.All: %v", err)
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			r.l.Errorf(ctx, "internal.comments.repository.postgres.Get.errChan: %v", err)
			return nil, paginator.Paginator{}, err
		}
	}
	dbComments := util.DerefSlice(cs)
	comments := make([]models.Comment, len(dbComments))
	for i, comment := range dbComments {
		comments[i] = models.NewComment(comment)
	}

	return comments, paginator.Paginator{
		Total:       total,
		Count:       int64(len(comments)),
		PerPage:     opts.PagQuery.Limit,
		CurrentPage: opts.PagQuery.Page,
	}, nil
}

func (r implRepository) Create(ctx context.Context, sc models.Scope, opts repository.CreateOptions) (models.Comment, error) {
	m := r.buildModel(ctx, sc, opts)

	err := m.Insert(ctx, r.database, boil.Infer())
	if err != nil {
		r.l.Errorf(ctx, "internal.comments.repository.postgres.Create.Insert: %v", err)
		return models.Comment{}, err
	}

	return models.NewComment(m), nil
}

func (r implRepository) Update(ctx context.Context, sc models.Scope, opts repository.UpdateOptions) (models.Comment, error) {
	c, col, err := r.buildUpdateModel(ctx, opts)
	if err != nil {
		r.l.Errorf(ctx, "internal.comments.repository.postgres.Update.buildUpdateModel: %v", err)
		return models.Comment{}, err
	}

	_, err = c.Update(ctx, r.database, boil.Whitelist(col...))
	if err != nil {
		r.l.Errorf(ctx, "internal.comments.repository.postgres.Update.Update: %v", err)
		return models.Comment{}, err
	}

	return models.NewComment(c), nil
}

func (r implRepository) Detail(ctx context.Context, sc models.Scope, ID string) (models.Comment, error) {
	qr, err := r.buildDetailQuery(ctx, ID)
	if err != nil {
		r.l.Errorf(ctx, "internal.comments.repository.postgres.Detail.buildDetailQuery: %v", err)
		return models.Comment{}, err
	}

	comment, err := dbmodels.Comments(qr...).One(ctx, r.database)
	if err != nil {
		if err == sql.ErrNoRows {
			r.l.Errorf(ctx, "internal.comments.repository.postgres.Detail.One.NoRows: %v", err)
			return models.Comment{}, repository.ErrNotFound
		}
		r.l.Errorf(ctx, "internal.comments.repository.postgres.Detail.One: %v", err)
		return models.Comment{}, err
	}

	return models.NewComment(*comment), nil
}

func (r implRepository) Delete(ctx context.Context, sc models.Scope, IDs []string) error {
	qr, err := r.buildDeleteQuery(ctx, IDs)
	if err != nil {
		r.l.Errorf(ctx, "internal.comments.repository.postgres.Delete.buildDeleteQuery: %v", err)
		return err
	}

	_, err = dbmodels.Comments(qr...).DeleteAll(ctx, r.database, true)
	if err != nil {
		r.l.Errorf(ctx, "internal.comments.repository.postgres.Delete.DeleteAll: %v", err)
		return err
	}

	return nil
}

func (r implRepository) GetByCard(ctx context.Context, sc models.Scope, cardID string) ([]models.Comment, error) {
	qr, err := r.buildGetByCardQuery(ctx, cardID)
	if err != nil {
		r.l.Errorf(ctx, "internal.comments.repository.postgres.GetByCard.buildGetByCardQuery: %v", err)
		return nil, err
	}

	cs, err := dbmodels.Comments(qr...).All(ctx, r.database)
	if err != nil {
		r.l.Errorf(ctx, "internal.comments.repository.postgres.GetByCard.All: %v", err)
		return nil, err
	}

	dbComments := util.DerefSlice(cs)
	comments := make([]models.Comment, len(dbComments))
	for i, comment := range dbComments {
		comments[i] = models.NewComment(comment)
	}

	return comments, nil
}
