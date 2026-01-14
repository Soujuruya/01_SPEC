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

// CreateIncident godoc
// @Summary Create Incidents
// @Description Создаёт новый инцидент с указанным заголовком, координатами и радиусом
// @Tags incident
// @Accept json
// @Produce json
// @Param request body incident.CreateIncidentRequest true "Incident creation data"
// @Success 201 {object} incident.IncidentResponse
// @Failure 400 {object} httphelper.APIResponse
// @Failure 500 {object} httphelper.APIResponse
// @Router /incidents [post]
func (h *IncidentHandler) CreateIncident(w http.ResponseWriter, r *http.Request) {
	var incidentDTO CreateIncidentRequest

	if err := json.NewDecoder(r.Body).Decode(&incidentDTO); err != nil {
		h.lg.Error("CreateIncident: failed to decode request", "error", err)
		httphelper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := ValidateCreateIncident(&incidentDTO); err != nil {
		h.lg.Error("CreateIncident: validation failed", "error", err)
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
		h.lg.Error("CreateIncident: failed to create incident", "error", err)
		httphelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.lg.Info("CreateIncident: incident created", "incident_id", inc.ID, "title", inc.Title)
	httphelper.WriteJSON(w, IncidentToResponse(inc), http.StatusCreated)
}

// GetIncident godoc
// @Summary Get Incident by ID
// @Description Получает инцидент по UUID
// @Tags incident
// @Produce json
// @Param id path string true "Incident UUID"
// @Success 200 {object} incident.IncidentResponse
// @Failure 400 {object} httphelper.APIResponse "Invalid UUID"
// @Failure 404 {object} httphelper.APIResponse "Incident not found"
// @Failure 500 {object} httphelper.APIResponse "Internal server error"
// @Router /incidents/{id} [get]
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

// GetListIncidents godoc
// @Summary Get All Incidents
// @Description Получает все активные и неактивные инциденты с поддержкой пагинации
// @Tags incident
// @Produce json
// @Param limit query int false "Limit for pagination" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} incident.IncidentListResponse
// @Failure 400 {object} httphelper.APIResponse "Invalid limit/offset"
// @Failure 500 {object} httphelper.APIResponse "Internal server error"
// @Router /incidents [get]
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

// GetActiveIncidents godoc
// @Summary Get All Active Incidents
// @Description Получает все активные инциденты
// @Tags incident
// @Produce json
// @Success 200 {object} incident.IncidentListResponse
// @Failure 500 {object} httphelper.APIResponse "Internal server error"
// @Router /incidents/active [get]
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

// DeactivateIncident godoc
// @Summary Deactivate Incident by ID
// @Description Деактивирует инцидент, не удаляя его полностью
// @Tags incident
// @Produce json
// @Param id path string true "Incident UUID"
// @Success 204 "No Content"
// @Failure 400 {object} httphelper.APIResponse "Invalid UUID"
// @Failure 404 {object} httphelper.APIResponse "Incident not found"
// @Failure 500 {object} httphelper.APIResponse "Internal server error"
// @Router /incidents/{id} [delete]
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

// UpdateIncident godoc
// @Summary Update Incident by ID
// @Description Обновляет данные об инциденте (частичное обновление)
// @Tags incident
// @Accept json
// @Produce json
// @Param id path string true "Incident UUID"
// @Param request body incident.UpdateIncidentRequest true "Incident update data"
// @Success 200 {object} incident.IncidentResponse
// @Failure 400 {object} httphelper.APIResponse "Invalid UUID or request body"
// @Failure 404 {object} httphelper.APIResponse "Incident not found"
// @Failure 500 {object} httphelper.APIResponse "Internal server error"
// @Router /incidents/{id} [patch]
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
