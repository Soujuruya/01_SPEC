package location

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LocationRepository interface {
	Save(ctx context.Context, loc *Location) error
	ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*Location, error)
	CountUniqueUsers(ctx context.Context, since time.Time) (int, error)
}
