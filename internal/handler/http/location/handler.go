package location

import (
	"encoding/json"
	"net/http"

	"github.com/Soujuruya/01_SPEC/internal/pkg/httphelper"
	"github.com/Soujuruya/01_SPEC/internal/usecase"
)

type LocationHandler struct {
	Service *usecase.LocationService
}

// CheckLocation возвращает последние локации пользователя
func (h *LocationHandler) CheckLocation(w http.ResponseWriter, r *http.Request) {
	var req CheckLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	loc, err := h.Service.CheckLocation(r.Context(), req.UserID, req.Lat, req.Lng)
	if err != nil {
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	httphelper.WriteJSON(w, LocationToResponse(loc), http.StatusOK)
}
