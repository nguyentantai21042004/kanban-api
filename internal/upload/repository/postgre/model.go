package postgre

import (
	"context"

	"github.com/volatiletech/null/v8"
	"gitlab.com/tantai-smap/authenticate-api/internal/dbmodels"
	"gitlab.com/tantai-smap/authenticate-api/internal/upload"
	"gitlab.com/tantai-smap/authenticate-api/pkg/postgre"
)

func (r implRepository) buildModel(ctx context.Context, opts upload.CreateOptions) dbmodels.Upload {
	upload := dbmodels.Upload{}

	if opts.Name != "" {
		upload.Name = opts.Name
	}

	if opts.Path != "" {
		upload.Path = opts.Path
	}

	if opts.Source != "" {
		upload.Source = opts.Source
	}

	if opts.FromLocation != "" {
		upload.FromLocation = opts.FromLocation
	}

	if opts.PublicID != "" {
		upload.PublicID = null.NewString(opts.PublicID, true)
	}

	if opts.CreatedUserID != "" {
		if err := postgre.IsUUID(opts.CreatedUserID); err != nil {
			r.l.Errorf(ctx, "internal.upload.repository.postgre.buildModel.IsUUID: %v", err)
			return dbmodels.Upload{}
		}
		upload.CreatedUserID = opts.CreatedUserID
	}

	upload.CreatedAt = null.NewTime(r.clock(), true)
	upload.UpdatedAt = null.NewTime(r.clock(), true)
	return upload
}
