package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `env-required:"true" env:"ENV"`

	DB    DBConfig    `env-required:"true" env-prefix:"DATABASE_"`
	Redis RedisConfig `env-required:"true" env-prefix:"REDIS_"`

	RetryLimit int           `env:"RETRY_LIMIT" env-default:"5"`
	RetryDelay time.Duration `env:"RETRY_DELAY" env-default:"5s"`

	WebhookURL             string `env-required:"true" env:"WEBHOOK_URL"`
	StatsTimeWindowMinutes int    `env:"STATS_TIME_WINDOW_MINUTES" env-default:"5"`

	HTTPPort      int           `env-required:"true" env:"HTTP_PORT"`
	HandleTimeout time.Duration `env-required:"true" env:"HANDLE_TIMEOUT"`
	CacheTTL      time.Duration `env-required:"true" env:"CACHE_TTL"`
}

type DBConfig struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME"`
}
type RedisConfig struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	Password string `env:"PASSWORD"`
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}

func Load(path string) (*Config, error) {
	var cfg Config

	if path != "" {
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			return nil, err
		}
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
