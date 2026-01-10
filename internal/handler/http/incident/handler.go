package incident

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Soujuruya/01_SPEC/internal/domain/incident"
	"github.com/Soujuruya/01_SPEC/internal/pkg/errs"
	"github.com/Soujuruya/01_SPEC/internal/pkg/httphelper"
	"github.com/Soujuruya/01_SPEC/internal/usecase"
)

type IncidentHandler struct {
	Service *usecase.IncidentService
}

func NewIncidentHandler(service *usecase.IncidentService) *IncidentHandler {
	return &IncidentHandler{
		Service: service,
	}
}

func (h *IncidentHandler) CreateIncidents(w http.ResponseWriter, r *http.Request) {
	var incidentDTO CreateIncidentRequest

	err := json.NewDecoder(r.Body).Decode(&incidentDTO)
	if err != nil {
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := ValidateCreateIncident(&incidentDTO); err != nil {
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	incident := incident.NewIncident(incidentDTO.Title, incidentDTO.Lat, incidentDTO.Lng, incidentDTO.Radius, true)
	err = h.Service.CreateIncident(r.Context(), incident)
	if err != nil {
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	incidentResDTO := IncidentToResponse(incident)

	httphelper.WriteJSON(w, incidentResDTO, http.StatusCreated)
}

func (h *IncidentHandler) GetIncident(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.ParseUUIDFromPath(r, "/api/v1/incidents/")
	if err != nil {
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	inc, err := h.Service.GetIncident(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			httphelper.WriteError(w, err, http.StatusNotFound)
			return
		}
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	incResponse := IncidentToResponse(inc)
	httphelper.WriteJSON(w, incResponse, http.StatusOK)
}

func (h *IncidentHandler) GetListIncidents(w http.ResponseWriter, r *http.Request) {
	var err error
	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			httphelper.WriteError(w, fmt.Errorf("invalid limit format"), http.StatusBadRequest)
			return
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			httphelper.WriteError(w, fmt.Errorf("invalid offset format"), http.StatusBadRequest)
			return
		}
	}

	if err := ValidateLimitOffset(limit, offset); err != nil {
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	incs, total, err := h.Service.ListIncidents(r.Context(), offset, limit)
	if err != nil {
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	incsResponse := IncidentsToListResponse(incs, offset, limit, total)

	httphelper.WriteJSON(w, incsResponse, http.StatusOK)
}

func (h *IncidentHandler) GetActiveIncidents(w http.ResponseWriter, r *http.Request) {
	activeIncs, err := h.Service.GetActiveIncidents(r.Context())
	if err != nil {
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	total, err := h.Service.CountActiveIncidents(r.Context())
	if err != nil {
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	activeIncsResponse := IncidentsToListResponse(activeIncs, 0, len(activeIncs), total)
	httphelper.WriteJSON(w, activeIncsResponse, http.StatusOK)
}

func (h *IncidentHandler) DeactivateIncident(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.ParseUUIDFromPath(r, "/api/v1/incidents/")
	if err != nil {
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	err = h.Service.DeactivateIncident(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			httphelper.WriteError(w, err, http.StatusNotFound)
			return
		}
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *IncidentHandler) UpdateIncident(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.ParseUUIDFromPath(r, "/api/v1/incidents/")
	if err != nil {
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}
	var incidentDTO UpdateIncidentRequest
	err = json.NewDecoder(r.Body).Decode(&incidentDTO)
	if err != nil {
		httphelper.WriteError(w, fmt.Errorf("invalid id format"), http.StatusBadRequest)
		return
	}
	if err := ValidateUpdateIncident(&incidentDTO); err != nil {
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	existing, err := h.Service.GetIncident(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			httphelper.WriteError(w, err, http.StatusNotFound)
			return
		}
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	if incidentDTO.Title != nil {
		existing.Title = *incidentDTO.Title
	}
	if incidentDTO.Lat != nil {
		existing.Lat = *incidentDTO.Lat
	}
	if incidentDTO.Lng != nil {
		existing.Lng = *incidentDTO.Lng
	}
	if incidentDTO.Radius != nil {
		existing.Radius = *incidentDTO.Radius
	}
	if incidentDTO.IsActive != nil {
		existing.IsActive = *incidentDTO.IsActive
	}
	existing.UpdatedAt = time.Now()

	err = h.Service.UpdateIncident(r.Context(), existing)
	if err != nil {
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	httphelper.WriteJSON(w, IncidentToResponse(existing), http.StatusOK)
}
