package location

import (
	"encoding/json"
	"net/http"

	"github.com/Soujuruya/01_SPEC/internal/pkg/httphelper"
	"github.com/Soujuruya/01_SPEC/internal/pkg/logger"
	"github.com/Soujuruya/01_SPEC/internal/usecase"
)

type LocationHandler struct {
	Service *usecase.LocationService
	lg      *logger.Logger
}

func NewLocationHandler(service *usecase.LocationService, lg *logger.Logger) *LocationHandler {
	return &LocationHandler{
		Service: service,
		lg:      lg,
	}
}

// CheckLocation godoc
// @Summary Check user location incidents
// @Description Возвращает локации инцидентов, в которые попал пользователь
// @Tags location
// @Accept json
// @Produce json
// @Param request body location.CheckLocationRequest true "User location data"
// @Success 200 {object} location.LocationResponse
// @Failure 400 {object} httphelper.APIResponse
// @Failure 500 {object} httphelper.APIResponse
// @Router /location/check [post]
func (h *LocationHandler) CheckLocation(w http.ResponseWriter, r *http.Request) {
	var req CheckLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.lg.Error("LocationHandler.CheckLocation: failed to decode request", "error", err)
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	loc, err := h.Service.CheckLocation(r.Context(), req.UserID, req.Lat, req.Lng)
	if err != nil {
		h.lg.Error("LocationHandler.CheckLocation: service returned error", "error", err, "user_id", req.UserID)
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.lg.Debug("LocationHandler.CheckLocation: location returned", "user_id", req.UserID, "location_id", loc.ID)
	httphelper.WriteJSON(w, LocationToResponse(loc), http.StatusOK)
}
