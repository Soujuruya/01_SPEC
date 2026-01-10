package stats

import (
	"net/http"

	"github.com/Soujuruya/01_SPEC/internal/config"
	"github.com/Soujuruya/01_SPEC/internal/pkg/httphelper"
	"github.com/Soujuruya/01_SPEC/internal/usecase"
)

type StatsHandler struct {
	Service *usecase.StatsService
	cfg     *config.Config
}

func NewStatsHandler(service *usecase.StatsService, cfg *config.Config) *StatsHandler {
	return &StatsHandler{
		Service: service,
		cfg:     cfg,
	}
}

// GetIncidentsStats возвращает кол-во уникальных пользователей за последние N минут
func (h *StatsHandler) GetIncidentsStats(w http.ResponseWriter, r *http.Request) {
	count, err := h.Service.GetUserCount(r.Context(), h.cfg.StatsTimeWindowMinutes)
	if err != nil {
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	resp := map[string]int{"user_count": count}
	httphelper.WriteJSON(w, resp, http.StatusOK)
}
