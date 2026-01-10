package health

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Soujuruya/01_SPEC/internal/pkg/logger"
	"go.uber.org/zap"
)

type HealthHandler struct {
	lg *logger.Logger
}

func NewHealthHandler(lg *logger.Logger) *HealthHandler {
	return &HealthHandler{lg: lg}
}

// HealthCheck возвращает базовый статус сервиса
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.lg.Debug("HealthHandler.HealthCheck: check received", zap.String("remote_addr", r.RemoteAddr))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.lg.Error("HealthHandler.HealthCheck: failed to write response",
			zap.Error(err))
	}
}
