package postgre

import (
	"context"
	"sync"

	"gitlab.com/tantai-smap/authenticate-api/internal/dbmodels"
	"gitlab.com/tantai-smap/authenticate-api/internal/models"
	"gitlab.com/tantai-smap/authenticate-api/internal/role"
	"gitlab.com/tantai-smap/authenticate-api/pkg/paginator"
	"gitlab.com/tantai-smap/authenticate-api/pkg/scope"
)

func (r implRepository) GetOne(ctx context.Context, sc scope.Scope, opts role.GetOneOptions) (models.Role, error) {
	qr, err := r.buildGetQuery(ctx, role.GetOptions{
		Filter: role.Filter(opts.Filter),
	})
	if err != nil {
		r.l.Errorf(ctx, "role.repository.postgre.GetOne.buildGetQuery: %v", err)
		return models.Role{}, err
	}

	dbRole, err := dbmodels.Roles(qr...).One(ctx, r.database)
	if err != nil {
		r.l.Errorf(ctx, "role.repository.postgre.GetOne.One: %v", err)
		return models.Role{}, err
	}

	return models.NewRole(*dbRole), nil
}

func (r implRepository) Detail(ctx context.Context, sc scope.Scope, ID string) (models.Role, error) {
	qr, err := r.buildDetailQuery(ctx, ID)
	if err != nil {
		r.l.Errorf(ctx, "role.repository.postgre.Detail.buildDetailQuery: %v", err)
		return models.Role{}, err
	}

	dbRole, err := dbmodels.Roles(qr...).One(ctx, r.database)
	if err != nil {
		r.l.Errorf(ctx, "role.repository.postgre.Detail.One: %v", err)
		return models.Role{}, err
	}

	return models.NewRole(*dbRole), nil
}

func (r implRepository) Get(ctx context.Context, sc scope.Scope, opts role.GetOptions) ([]models.Role, paginator.Paginator, error) {
	qr, err := r.buildGetQuery(ctx, opts)
	if err != nil {
		r.l.Errorf(ctx, "role.repository.postgre.Get.buildGetQuery: %v", err)
		return nil, paginator.Paginator{}, err
	}

	var (
		wg    sync.WaitGroup
		wgErr error
		roles []models.Role
		total int64
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := dbmodels.Roles(qr...).Bind(ctx, r.database, &roles)
		if err != nil {
			r.l.Errorf(ctx, "role.repository.postgre.Get.Bind: %v", err)
			wgErr = err
			return
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		total, err = dbmodels.Roles(qr...).Count(ctx, r.database)
		if err != nil {
			r.l.Errorf(ctx, "user.repository.postgre.Get.Count: %v", err)
			wgErr = err
			return
		}
	}()

	wg.Wait()
	if wgErr != nil {
		return nil, paginator.Paginator{}, wgErr
	}

	return roles, paginator.Paginator{
		Total:       total,
		Count:       int64(len(roles)),
		PerPage:     opts.PagQuery.Limit,
		CurrentPage: opts.PagQuery.Page,
	}, nil
}

func (r implRepository) List(ctx context.Context, sc scope.Scope, opts role.ListOptions) ([]models.Role, error) {
	qr, err := r.buildGetQuery(ctx, role.GetOptions{
		Filter: role.Filter(opts.Filter),
	})
	if err != nil {
		r.l.Errorf(ctx, "role.repository.postgre.List.buildGetQuery: %v", err)
		return nil, err
	}

	var roles []models.Role
	err = dbmodels.Roles(qr...).Bind(ctx, r.database, &roles)
	if err != nil {
		r.l.Errorf(ctx, "role.repository.postgre.List.Bind: %v", err)
		return nil, err
	}

	return roles, nil
}
