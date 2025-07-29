package postgre

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"gitlab.com/tantai-smap/authenticate-api/internal/dbmodels"
	"gitlab.com/tantai-smap/authenticate-api/internal/models"
	"gitlab.com/tantai-smap/authenticate-api/internal/upload"
	"gitlab.com/tantai-smap/authenticate-api/pkg/scope"
)

func (r implRepository) Detail(ctx context.Context, sc scope.Scope, ID string) (models.Upload, error) {
	qr, err := r.buildDetailQuery(ctx, ID)
	if err != nil {
		r.l.Errorf(ctx, "internal.upload.repository.postgre.Detail.buildDetailQuery: %v", err)
		return models.Upload{}, err
	}

	dbUpload, err := dbmodels.Uploads(qr...).One(ctx, r.database)
	if err != nil {
		r.l.Errorf(ctx, "internal.upload.repository.postgre.Detail.One: %v", err)
		return models.Upload{}, err
	}

	return models.NewUpload(*dbUpload), nil
}

func (r implRepository) Create(ctx context.Context, sc scope.Scope, opts upload.CreateOptions) (models.Upload, error) {
	dbUpload := r.buildModel(ctx, opts)

	err := dbUpload.Insert(ctx, r.database, boil.Infer())
	if err != nil {
		r.l.Errorf(ctx, "internal.upload.repository.postgre.Create.Insert: %v", err)
		return models.Upload{}, err
	}

	return models.NewUpload(dbUpload), nil
}
