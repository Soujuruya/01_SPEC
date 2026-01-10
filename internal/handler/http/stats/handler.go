package stats

import (
	"net/http"

	"github.com/Soujuruya/01_SPEC/internal/config"
	"github.com/Soujuruya/01_SPEC/internal/pkg/httphelper"
	"github.com/Soujuruya/01_SPEC/internal/pkg/logger"
	"github.com/Soujuruya/01_SPEC/internal/usecase"
)

type StatsHandler struct {
	Service *usecase.StatsService
	cfg     *config.Config
	lg      *logger.Logger
}

func NewStatsHandler(service *usecase.StatsService, cfg *config.Config, lg *logger.Logger) *StatsHandler {
	return &StatsHandler{
		Service: service,
		cfg:     cfg,
		lg:      lg,
	}
}

// GetIncidentsStats возвращает кол-во уникальных пользователей за последние N минут
func (h *StatsHandler) GetIncidentsStats(w http.ResponseWriter, r *http.Request) {
	count, err := h.Service.GetUserCount(r.Context(), h.cfg.StatsTimeWindowMinutes)
	if err != nil {
		h.lg.Error("StatsHandler.GetIncidentsStats: failed to get user count", "error", err)
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.lg.Debug("StatsHandler.GetIncidentsStats: user count returned", "count", count, "minutes", h.cfg.StatsTimeWindowMinutes)
	resp := map[string]int{"user_count": count}
	httphelper.WriteJSON(w, resp, http.StatusOK)
}
