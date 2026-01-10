package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Soujuruya/01_SPEC/internal/domain/location"
	"github.com/Soujuruya/01_SPEC/internal/integration"
)

type WebhookQueue struct {
	rdb *redis.Client
	key string
}

func NewWebhookQueue(rdb *redis.Client, key string) *WebhookQueue {
	return &WebhookQueue{
		rdb: rdb,
		key: key,
	}
}

func (q *WebhookQueue) Enqueue(ctx context.Context, loc *location.Location) error {
	payload := integration.WebhookPayload{
		UserID:      loc.UserID,
		Lat:         loc.Lat,
		Lng:         loc.Lng,
		IncidentIDs: loc.IncidentIDs,
		Timestamp:   time.Now().Unix(),
		Retry:       0,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return q.rdb.LPush(ctx, q.key, data).Err()
}
