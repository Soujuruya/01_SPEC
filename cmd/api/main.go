package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Soujuruya/01_SPEC/internal/config"
	"github.com/Soujuruya/01_SPEC/internal/handler/http/health"
	"github.com/Soujuruya/01_SPEC/internal/handler/http/incident"
	"github.com/Soujuruya/01_SPEC/internal/handler/http/location"
	"github.com/Soujuruya/01_SPEC/internal/handler/http/stats"
	"github.com/Soujuruya/01_SPEC/internal/integration"
	redispkg "github.com/Soujuruya/01_SPEC/internal/pkg/redis"
	"github.com/Soujuruya/01_SPEC/internal/repository/postgres"
	"github.com/Soujuruya/01_SPEC/internal/repository/redis"
	"github.com/Soujuruya/01_SPEC/internal/server"
	"github.com/Soujuruya/01_SPEC/internal/usecase"
	"github.com/Soujuruya/01_SPEC/internal/worker"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Конфиг
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Postgres
	dbURL := cfg.DB.DSN()
	pgxPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}
	defer pgxPool.Close()

	rdb := redispkg.NewClient(&cfg.Redis)
	defer func() {
		if err := rdb.Close(); err != nil {
			fmt.Println("error closing Redis:", err)
		}
	}()

	// 4. Репозитории
	incidentRepo := postgres.NewIncidentRepo(pgxPool)
	locationRepo := postgres.NewLocationRepo(pgxPool)

	// 5. Кэш и очередь
	incidentCache := redis.NewIncidentCache(rdb, "active_incidents", cfg.CacheTTL)
	webhookQueue := redis.NewWebhookQueue(rdb, "webhook_queue")

	// 6. Сервисы
	incidentService := usecase.NewIncidentService(incidentRepo, incidentCache)
	locationService := usecase.NewLocationService(locationRepo, incidentRepo, webhookQueue)
	statsService := usecase.NewStatsService(locationRepo)
	fmt.Printf("%+v\n", incidentService)

	// 6.1. Воркер для вебхуков
	webhookClient := integration.NewWebhookClient(cfg.WebhookURL, cfg.HandleTimeout)
	worker := worker.NewWebhookWorker(rdb, "webhook_queue", webhookClient, cfg.RetryLimit, cfg.RetryDelay)
	go worker.Run(context.Background())

	// 7. Хендлеры
	healthHandler := health.NewHealthHandler()
	incidentHandler := incident.NewIncidentHandler(incidentService)
	locationHandler := location.NewLocationHandler(locationService)
	statsHandler := stats.NewStatsHandler(statsService, cfg)
	fmt.Printf("%+v\n", incidentHandler.Service)

	srv := server.NewServer(cfg, healthHandler, incidentHandler, locationHandler, statsHandler)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HandleTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	log.Println("server stopped gracefully")
}
