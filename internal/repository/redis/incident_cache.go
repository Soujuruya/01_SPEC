package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Soujuruya/01_SPEC/internal/domain/incident"
	"github.com/Soujuruya/01_SPEC/internal/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type IncidentCache struct {
	rdb *redis.Client
	key string
	ttl time.Duration
	lg  *logger.Logger
}

func NewIncidentCache(
	rdb *redis.Client,
	key string,
	ttl time.Duration,
	lg *logger.Logger,
) *IncidentCache {
	return &IncidentCache{
		rdb: rdb,
		key: key,
		ttl: ttl,
		lg:  lg,
	}
}

func (c *IncidentCache) GetActive(ctx context.Context) ([]*incident.Incident, error) {
	data, err := c.rdb.Get(ctx, c.key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.lg.Debug("GetActive: cache miss", "key", c.key)
			return nil, err
		}
		c.lg.Error("GetActive: failed to get from Redis", "key", c.key, "error", err)
		return nil, err
	}

	var incs []*incident.Incident
	if err := json.Unmarshal(data, &incs); err != nil {
		c.lg.Error("GetActive: failed to unmarshal cached data", "key", c.key, "error", err)
		return nil, err
	}

	c.lg.Debug("GetActive: cache hit", "key", c.key, "count", len(incs))
	return incs, nil
}

func (c *IncidentCache) SetActive(ctx context.Context, incs []*incident.Incident) error {
	bytes, err := json.Marshal(incs)
	if err != nil {
		c.lg.Error("SetActive: failed to marshal incidents", "key", c.key, "error", err)
		return err
	}

	if err := c.rdb.Set(ctx, c.key, bytes, c.ttl).Err(); err != nil {
		c.lg.Error("SetActive: failed to set Redis key", "key", c.key, "error", err)
		return err
	}

	c.lg.Debug("SetActive: cache updated", "key", c.key, "count", len(incs))
	return nil
}

func (c *IncidentCache) InvalidateActive(ctx context.Context) error {
	err := c.rdb.Del(ctx, c.key).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		c.lg.Error("InvalidateActive: failed to delete Redis key", "key", c.key, "error", err)
		return err
	}

	c.lg.Debug("InvalidateActive: cache invalidated", "key", c.key)
	return nil
}
