package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Soujuruya/01_SPEC/internal/config"
	"github.com/Soujuruya/01_SPEC/internal/handler/http/health"
	"github.com/Soujuruya/01_SPEC/internal/handler/http/incident"
	"github.com/Soujuruya/01_SPEC/internal/handler/http/location"
	"github.com/Soujuruya/01_SPEC/internal/handler/http/middleware"
	"github.com/Soujuruya/01_SPEC/internal/handler/http/stats"
	"github.com/Soujuruya/01_SPEC/internal/integration"
	"github.com/Soujuruya/01_SPEC/internal/pkg/logger"
	redispkg "github.com/Soujuruya/01_SPEC/internal/pkg/redis"
	"github.com/Soujuruya/01_SPEC/internal/repository/postgres"
	"github.com/Soujuruya/01_SPEC/internal/repository/redis"
	"github.com/Soujuruya/01_SPEC/internal/server"
	"github.com/Soujuruya/01_SPEC/internal/usecase"
	"github.com/Soujuruya/01_SPEC/internal/worker"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Конфиг
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	//Логгер
	lg := logger.New(cfg.Environment)
	defer func() { _ = lg.Sync() }()

	//  Postgres
	dbURL := cfg.DB.DSN()
	pgxPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("failed to connect postgres", "error", err)
	}
	defer pgxPool.Close()
	lg.Info("postgres connected")

	// Redis
	rdb := redispkg.NewClient(&cfg.Redis)
	defer func() {
		if err := rdb.Close(); err != nil {
			lg.Error("failed to close Redis", "error", err)
		}
	}()
	lg.Info("redis connected")

	//  Репозитории
	incidentRepo := postgres.NewIncidentRepo(pgxPool, lg)
	locationRepo := postgres.NewLocationRepo(pgxPool, lg)

	// Кэш и очередь
	incidentCache := redis.NewIncidentCache(rdb, "active_incidents", cfg.CacheTTL, lg)
	webhookQueue := redis.NewWebhookQueue(rdb, "webhook_queue", lg)

	//  Сервисы
	incidentService := usecase.NewIncidentService(incidentRepo, incidentCache, lg)
	locationService := usecase.NewLocationService(locationRepo, incidentRepo, webhookQueue, lg)
	statsService := usecase.NewStatsService(locationRepo, lg)

	// Воркер для вебхуков
	webhookClient := integration.NewWebhookClient(cfg.WebhookURL, cfg.HandleTimeout, lg)
	worker := worker.NewWebhookWorker(rdb, "webhook_queue", webhookClient, cfg.RetryLimit, cfg.RetryDelay, lg)
	go worker.Run(context.Background(), lg)

	// Хендлеры
	healthHandler := health.NewHealthHandler(lg)
	incidentHandler := incident.NewIncidentHandler(incidentService, lg)
	locationHandler := location.NewLocationHandler(locationService, lg)
	statsHandler := stats.NewStatsHandler(statsService, cfg, lg)

	//  HTTP Server
	srv := server.NewServer(cfg,
		healthHandler,
		incidentHandler,
		locationHandler,
		statsHandler,
		middleware.Logger(lg), //  middleware логирования
	)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal("server failed", "error", err)
		}
	}()

	//Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HandleTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown failed", "error", err)
	}

	lg.Info("server stopped gracefully")
}
