package location

import (
	"context"

	"github.com/google/uuid"
)

type LocationRepository interface {
	Save(ctx context.Context, loc *Location) error
	ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*Location, error)
}
