package location

import "github.com/Soujuruya/01_SPEC/internal/domain/location"

func LocationToResponse(loc *location.Location) *LocationResponse {
	return &LocationResponse{
		ID:          loc.ID,
		UserID:      loc.UserID,
		Lat:         loc.Lat,
		Lng:         loc.Lng,
		Timestamp:   loc.Timestamp,
		IsCheck:     loc.IsCheck,
		IncidentIDs: loc.IncidentIDs,
	}
}
