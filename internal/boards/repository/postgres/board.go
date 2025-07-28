package postgres

import (
	"context"
	"sync"

	"github.com/aarondl/sqlboiler/v4/boil"
	"gitlab.com/tantai-kanban/kanban-api/internal/boards/repository"
	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
	"gitlab.com/tantai-kanban/kanban-api/pkg/util"
)

func (r implRepository) Get(ctx context.Context, sc models.Scope, opts repository.GetOptions) ([]models.Board, paginator.Paginator, error) {
	qr, err := r.buildGetQuery(ctx, opts.Filter)
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Get.buildGetQuery: %v", err)
		return nil, paginator.Paginator{}, err
	}

	var (
		total int64
		bs    dbmodels.BoardSlice
	)

	errChan := make(chan error, 2)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		var countErr error
		total, countErr = dbmodels.Boards(qr...).Count(ctx, r.database)
		if countErr != nil {
			r.l.Errorf(ctx, "internal.boards.repository.postgres.Get.Count: %v", countErr)
			errChan <- countErr
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bs, err = dbmodels.Boards(qr...).All(ctx, r.database)
		if err != nil {
			r.l.Errorf(ctx, "internal.boards.repository.postgres.Get.All: %v", err)
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			r.l.Errorf(ctx, "internal.boards.repository.postgres.Get.errChan: %v", err)
			return nil, paginator.Paginator{}, err
		}
	}
	dbBoards := util.DerefSlice(bs)
	boards := make([]models.Board, len(dbBoards))
	for i, board := range dbBoards {
		boards[i] = models.NewBoard(board)
	}

	return boards, paginator.Paginator{
		Total:       total,
		Count:       int64(len(boards)),
		PerPage:     opts.PagQuery.Limit,
		CurrentPage: opts.PagQuery.Page,
	}, nil
}

func (r implRepository) Create(ctx context.Context, sc models.Scope, opts repository.CreateOptions) (models.Board, error) {
	m := r.buildModel(ctx, opts)

	err := m.Insert(ctx, r.database, boil.Infer())
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Create.Insert: %v", err)
		return models.Board{}, err
	}

	return models.NewBoard(m), nil
}

func (r implRepository) Update(ctx context.Context, sc models.Scope, opts repository.UpdateOptions) (models.Board, error) {
	b, col, err := r.buildUpdateModel(ctx, opts)
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Update.buildUpdateModel: %v", err)
		return models.Board{}, err
	}

	_, err = b.Update(ctx, r.database, boil.Whitelist(col...))
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Update.Update: %v", err)
		return models.Board{}, err
	}

	return models.NewBoard(b), nil
}

func (r implRepository) Detail(ctx context.Context, sc models.Scope, ID string) (models.Board, error) {
	qr, err := r.buildDetailQuery(ctx, ID)
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Detail.buildDetailQuery: %v", err)
		return models.Board{}, err
	}

	board, err := dbmodels.Boards(qr...).One(ctx, r.database)
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Detail.One: %v", err)
		return models.Board{}, err
	}

	return models.NewBoard(*board), nil
}

func (r implRepository) Delete(ctx context.Context, sc models.Scope, IDs []string) error {
	qr, err := r.buildDeleteQuery(ctx, IDs)
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Delete.buildDeleteQuery: %v", err)
		return err
	}

	_, err = dbmodels.Boards(qr...).DeleteAll(ctx, r.database, true)
	if err != nil {
		r.l.Errorf(ctx, "internal.boards.repository.postgres.Delete.DeleteAll: %v", err)
		return err
	}

	return nil
}
