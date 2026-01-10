package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Soujuruya/01_SPEC/internal/domain/location"
	"github.com/Soujuruya/01_SPEC/internal/integration"
	"github.com/Soujuruya/01_SPEC/internal/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type WebhookQueue struct {
	rdb *redis.Client
	key string
	lg  *logger.Logger
}

func NewWebhookQueue(rdb *redis.Client, key string, lg *logger.Logger) *WebhookQueue {
	return &WebhookQueue{
		rdb: rdb,
		key: key,
		lg:  lg,
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
		q.lg.Error("WebhookQueue.Enqueue: failed to marshal payload", "user_id", loc.UserID, "error", err)
		return err
	}

	if err := q.rdb.LPush(ctx, q.key, data).Err(); err != nil {
		q.lg.Error("WebhookQueue.Enqueue: failed to push to Redis queue", "key", q.key, "user_id", loc.UserID, "error", err)
		return err
	}

	q.lg.Debug("WebhookQueue.Enqueue: task enqueued successfully", "key", q.key, "user_id", loc.UserID, "incident_count", len(loc.IncidentIDs))
	return nil
}
