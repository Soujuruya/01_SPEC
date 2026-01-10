package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Soujuruya/01_SPEC/internal/domain/incident"
	"github.com/redis/go-redis/v9"
)

type IncidentCache struct {
	rdb *redis.Client
	key string
	ttl time.Duration
}

func NewIncidentCache(
	rdb *redis.Client,
	key string,
	ttl time.Duration,
) *IncidentCache {
	return &IncidentCache{
		rdb: rdb,
		key: key,
		ttl: ttl,
	}
}

func (c *IncidentCache) GetActive(ctx context.Context) ([]*incident.Incident, error) {
	data, err := c.rdb.Get(ctx, c.key).Bytes()
	if err != nil {
		return nil, err
	}

	var incs []*incident.Incident
	if err := json.Unmarshal(data, &incs); err != nil {
		return nil, err
	}

	return incs, nil
}

func (c *IncidentCache) SetActive(ctx context.Context, incs []*incident.Incident) error {
	bytes, err := json.Marshal(incs)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, c.key, bytes, c.ttl).Err()
}

func (c *IncidentCache) InvalidateActive(ctx context.Context) error {
	err := c.rdb.Del(ctx, c.key).Err()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	return err
}
