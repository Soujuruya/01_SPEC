package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Soujuruya/01_SPEC/internal/config"
	"github.com/Soujuruya/01_SPEC/internal/handler/http/health"
	"github.com/Soujuruya/01_SPEC/internal/handler/http/incident"
	"github.com/Soujuruya/01_SPEC/internal/handler/http/location"
	"github.com/Soujuruya/01_SPEC/internal/handler/http/stats"
)

type Middleware func(http.Handler) http.Handler

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
}

func NewServer(cfg *config.Config,
	healthHandler *health.HealthHandler,
	incidentHandler *incident.IncidentHandler,
	locationHandler *location.LocationHandler,
	statsHandler *stats.StatsHandler,
	middlewares ...Middleware,
) *Server {

	mux := http.NewServeMux()

	// Health-check
	mux.HandleFunc("/api/v1/system/health", healthHandler.HealthCheck)

	// Список инцидентов, создание
	mux.HandleFunc("/api/v1/incidents", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			incidentHandler.GetListIncidents(w, r)
		case http.MethodPost:
			incidentHandler.CreateIncidents(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Конкретный инцидент (GET/PUT/DELETE)
	mux.HandleFunc("/api/v1/incidents/", func(w http.ResponseWriter, r *http.Request) {
		idPath := r.URL.Path[len("/api/v1/incidents/"):]
		if idPath == "" || idPath == "active" || idPath == "stats" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			incidentHandler.GetIncident(w, r)
		case http.MethodPut:
			incidentHandler.UpdateIncident(w, r)
		case http.MethodDelete:
			incidentHandler.DeactivateIncident(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Активные инциденты
	mux.HandleFunc("/api/v1/incidents/active", incidentHandler.GetActiveIncidents)

	// Проверка локации
	mux.HandleFunc("/api/v1/location/check", locationHandler.CheckLocation)

	// Статистика
	mux.HandleFunc("/api/v1/incidents/stats", statsHandler.GetIncidentsStats)

	var handler http.Handler = mux
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:      handler,
		ReadTimeout:  cfg.HandleTimeout,
		WriteTimeout: cfg.HandleTimeout,
		IdleTimeout:  cfg.HandleTimeout,
	}

	return &Server{
		httpServer: srv,
		cfg:        cfg,
	}
}

func (s *Server) Start() error {
	fmt.Printf("HTTP server listening on %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
