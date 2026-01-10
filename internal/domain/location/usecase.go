package location

import (
	"context"
)

type WebhookQueue interface {
	Enqueue(ctx context.Context, loc *Location) error
}
