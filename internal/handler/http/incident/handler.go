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
	"github.com/Soujuruya/01_SPEC/internal/pkg/logger"
	"github.com/Soujuruya/01_SPEC/internal/usecase"
)

type IncidentHandler struct {
	Service *usecase.IncidentService
	lg      *logger.Logger
}

func NewIncidentHandler(service *usecase.IncidentService, lg *logger.Logger) *IncidentHandler {
	return &IncidentHandler{
		Service: service,
		lg:      lg,
	}
}

func (h *IncidentHandler) CreateIncidents(w http.ResponseWriter, r *http.Request) {
	var incidentDTO CreateIncidentRequest

	if err := json.NewDecoder(r.Body).Decode(&incidentDTO); err != nil {
		h.lg.Error("CreateIncidents: failed to decode request", "error", err)
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := ValidateCreateIncident(&incidentDTO); err != nil {
		h.lg.Error("CreateIncidents: validation failed", "error", err)
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	inc := incident.NewIncident(
		incidentDTO.Title,
		incidentDTO.Lat,
		incidentDTO.Lng,
		incidentDTO.Radius,
		true,
	)

	if err := h.Service.CreateIncident(r.Context(), inc); err != nil {
		h.lg.Error("CreateIncidents: failed to create incident", "error", err)
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.lg.Info("CreateIncidents: incident created", "incident_id", inc.ID, "title", inc.Title)
	httphelper.WriteJSON(w, IncidentToResponse(inc), http.StatusCreated)
}

func (h *IncidentHandler) GetIncident(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.ParseUUIDFromPath(r, "/api/v1/incidents/")
	if err != nil {
		h.lg.Error("GetIncident: invalid UUID in path", "error", err)
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	inc, err := h.Service.GetIncident(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			h.lg.Warn("GetIncident: incident not found", "incident_id", id)
			httphelper.WriteError(w, err, http.StatusNotFound)
			return
		}
		h.lg.Error("GetIncident: failed to fetch incident", "incident_id", id, "error", err)
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.lg.Debug("GetIncident: success", "incident_id", id)
	httphelper.WriteJSON(w, IncidentToResponse(inc), http.StatusOK)
}

func (h *IncidentHandler) GetListIncidents(w http.ResponseWriter, r *http.Request) {
	limit := 10
	offset := 0
	var err error

	if l := r.URL.Query().Get("limit"); l != "" {
		limit, err = strconv.Atoi(l)
		if err != nil {
			h.lg.Error("GetListIncidents: invalid limit", "value", l, "error", err)
			httphelper.WriteError(w, fmt.Errorf("invalid limit format"), http.StatusBadRequest)
			return
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		offset, err = strconv.Atoi(o)
		if err != nil {
			h.lg.Error("GetListIncidents: invalid offset", "value", o, "error", err)
			httphelper.WriteError(w, fmt.Errorf("invalid offset format"), http.StatusBadRequest)
			return
		}
	}

	if err := ValidateLimitOffset(limit, offset); err != nil {
		h.lg.Error("GetListIncidents: invalid limit/offset", "limit", limit, "offset", offset, "error", err)
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	incs, total, err := h.Service.ListIncidents(r.Context(), offset, limit)
	if err != nil {
		h.lg.Error("GetListIncidents: failed to list incidents", "offset", offset, "limit", limit, "error", err)
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.lg.Debug("GetListIncidents: success", "returned", len(incs), "total", total)
	httphelper.WriteJSON(w, IncidentsToListResponse(incs, offset, limit, total), http.StatusOK)
}

func (h *IncidentHandler) GetActiveIncidents(w http.ResponseWriter, r *http.Request) {
	incs, err := h.Service.GetActiveIncidents(r.Context())
	if err != nil {
		h.lg.Error("GetActiveIncidents: failed to fetch active incidents", "error", err)
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	total, err := h.Service.CountActiveIncidents(r.Context())
	if err != nil {
		h.lg.Error("GetActiveIncidents: failed to count active incidents", "error", err)
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.lg.Debug("GetActiveIncidents: success", "count", len(incs), "total", total)
	httphelper.WriteJSON(w, IncidentsToListResponse(incs, 0, len(incs), total), http.StatusOK)
}

func (h *IncidentHandler) DeactivateIncident(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.ParseUUIDFromPath(r, "/api/v1/incidents/")
	if err != nil {
		h.lg.Error("DeactivateIncident: invalid UUID in path", "error", err)
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := h.Service.DeactivateIncident(r.Context(), id); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			h.lg.Warn("DeactivateIncident: incident not found", "incident_id", id)
			httphelper.WriteError(w, err, http.StatusNotFound)
			return
		}
		h.lg.Error("DeactivateIncident: failed to deactivate", "incident_id", id, "error", err)
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.lg.Info("DeactivateIncident: success", "incident_id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *IncidentHandler) UpdateIncident(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.ParseUUIDFromPath(r, "/api/v1/incidents/")
	if err != nil {
		h.lg.Error("UpdateIncident: invalid UUID in path", "error", err)
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	var incidentDTO UpdateIncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&incidentDTO); err != nil {
		h.lg.Error("UpdateIncident: failed to decode request body", "error", err)
		httphelper.WriteError(w, fmt.Errorf("invalid request body"), http.StatusBadRequest)
		return
	}

	if err := ValidateUpdateIncident(&incidentDTO); err != nil {
		h.lg.Error("UpdateIncident: validation failed", "error", err)
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	existing, err := h.Service.GetIncident(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			h.lg.Warn("UpdateIncident: incident not found", "incident_id", id)
			httphelper.WriteError(w, err, http.StatusNotFound)
			return
		}
		h.lg.Error("UpdateIncident: failed to get existing incident", "incident_id", id, "error", err)
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

	if err := h.Service.UpdateIncident(r.Context(), existing); err != nil {
		h.lg.Error("UpdateIncident: failed to update incident", "incident_id", id, "error", err)
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.lg.Info("UpdateIncident: success", "incident_id", id)
	httphelper.WriteJSON(w, IncidentToResponse(existing), http.StatusOK)
}
