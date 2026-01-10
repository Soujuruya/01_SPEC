package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Soujuruya/01_SPEC/internal/integration"
	"github.com/redis/go-redis/v9"
)

type WebhookWorker struct {
	rdb        *redis.Client
	queueKey   string
	client     *integration.WebhookClient
	retryMax   int
	retryDelay time.Duration
}

func NewWebhookWorker(
	rdb *redis.Client,
	queueKey string,
	client *integration.WebhookClient,
	retryMax int,
	retryDelay time.Duration,
) *WebhookWorker {
	return &WebhookWorker{
		rdb:        rdb,
		queueKey:   queueKey,
		client:     client,
		retryMax:   retryMax,
		retryDelay: retryDelay,
	}
}

func (w *WebhookWorker) Run(ctx context.Context) {
	for {
		res, err := w.rdb.BRPop(ctx, 0, w.queueKey).Result()
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			continue
		}

		var payload integration.WebhookPayload
		if err := json.Unmarshal([]byte(res[1]), &payload); err != nil {
			continue
		}

		err = w.client.Send(ctx, payload)
		if err == nil {
			continue
		}

		payload.Retry++
		if payload.Retry >= w.retryMax {
			continue
		}

		time.Sleep(w.retryDelay)

		data, _ := json.Marshal(payload)
		_ = w.rdb.LPush(ctx, w.queueKey, data).Err()
	}
}
