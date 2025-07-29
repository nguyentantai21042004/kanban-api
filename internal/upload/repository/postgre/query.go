package postgre

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gitlab.com/tantai-smap/authenticate-api/pkg/postgre"
)

func (r implRepository) buildDetailQuery(ctx context.Context, ID string) ([]qm.QueryMod, error) {
	qr := postgre.BuildQueryWithSoftDelete()

	if err := postgre.IsUUID(ID); err != nil {
		r.l.Errorf(ctx, "internal.upload.repository.postgre.buildDetailQuery.InvalidID: %v", err)
		return nil, err
	}
	qr = append(qr, qm.Where("id = ?", ID))

	return qr, nil
}
