package postgre

import (
	"github.com/aarondl/sqlboiler/queries/qm"
	"github.com/google/uuid"
)

// IsUUID checks if the given string is a valid UUID
func IsUUID(u string) error {
	_, err := uuid.Parse(u)
	if err != nil {
		return err
	}
	return nil
}

func NewUUID() string {
	return uuid.New().String()
}

func BuildQueryWithSoftDelete() []qm.QueryMod {
	return []qm.QueryMod{
		qm.Where("deleted_at IS NULL"),
	}
}
