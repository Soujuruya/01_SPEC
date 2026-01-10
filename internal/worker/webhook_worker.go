package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Soujuruya/01_SPEC/internal/integration"
	"github.com/Soujuruya/01_SPEC/internal/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type WebhookWorker struct {
	rdb        *redis.Client
	queueKey   string
	client     *integration.WebhookClient
	retryMax   int
	retryDelay time.Duration
	lg         *logger.Logger
}

func NewWebhookWorker(
	rdb *redis.Client,
	queueKey string,
	client *integration.WebhookClient,
	retryMax int,
	retryDelay time.Duration,
	lg *logger.Logger,
) *WebhookWorker {
	return &WebhookWorker{
		rdb:        rdb,
		queueKey:   queueKey,
		client:     client,
		retryMax:   retryMax,
		retryDelay: retryDelay,
		lg:         lg,
	}
}

func (w *WebhookWorker) Run(ctx context.Context, lg *logger.Logger) {
	for {
		res, err := w.rdb.BRPop(ctx, 0, w.queueKey).Result()
		if err != nil {
			if ctx.Err() != nil {
				w.lg.Info("WebhookWorker stopped due to context cancellation")
				return
			}
			w.lg.Error("WebhookWorker BRPop error", "error", err)
			continue
		}

		var payload integration.WebhookPayload
		if err := json.Unmarshal([]byte(res[1]), &payload); err != nil {
			w.lg.Error("Failed to unmarshal webhook payload", "error", err, "data", res[1])
			continue
		}

		err = w.client.Send(ctx, payload, lg)
		if err == nil {
			w.lg.Info("Webhook sent successfully", "user_id", payload.UserID, "incident_ids", payload.IncidentIDs)
			continue
		}

		payload.Retry++
		if payload.Retry >= w.retryMax {
			w.lg.Warn("Webhook retry limit reached, dropping payload", "user_id", payload.UserID, "incident_ids", payload.IncidentIDs, "retries", payload.Retry)
			continue
		}

		w.lg.Warn("Webhook send failed, retrying", "user_id", payload.UserID, "incident_ids", payload.IncidentIDs, "retry", payload.Retry, "error", err)
		time.Sleep(w.retryDelay)

		data, _ := json.Marshal(payload)
		if err := w.rdb.LPush(ctx, w.queueKey, data).Err(); err != nil {
			w.lg.Error("Failed to push webhook back to queue", "error", err)
		}
	}
}
