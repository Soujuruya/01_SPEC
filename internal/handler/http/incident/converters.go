package incident

import (
	"time"

	"github.com/Soujuruya/01_SPEC/internal/domain/incident"
)

func IncidentToResponse(inc *incident.Incident) IncidentResponse {
	return IncidentResponse{
		ID:        inc.ID.String(),
		Title:     inc.Title,
		Lat:       inc.Lat,
		Lng:       inc.Lng,
		IsActive:  inc.IsActive,
		CreatedAt: inc.CreatedAt.Format(time.RFC3339),
	}
}

func IncidentsToListResponse(incs []*incident.Incident, offset, limit, total int) IncidentListResponse {
	incsResponse := make([]IncidentResponse, len(incs))
	for i, inc := range incs {
		incsResponse[i] = IncidentToResponse(inc)
	}
	return IncidentListResponse{
		Incidents: incsResponse,
		Limit:     limit,
		Offset:    offset,
		Total:     total,
	}
}
