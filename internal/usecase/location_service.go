package usecase

import (
	"time"

	"github.com/Soujuruya/01_SPEC/internal/domain/location"
	"github.com/google/uuid"
)

type LocationService struct {
	Repo         location.LocationRepository
	Queue        location.WebhookQueue
	IncidentRepo incident.IncidentRepository
}

func (s *LocationService) CheckLocation(userID uuid.UUID, lat, lng float64) (*location.Location, error) {
	activeIncidents, _ := s.IncidentRepo.GetActiveIncidents()

	loc := &location.Location{
		ID:          uuid.New(),
		UserID:      userID,
		Lat:         lat,
		Lng:         lng,
		Timestamp:   time.Now(),
		IncidentIDs: []uuid.UUID{},
	}
	for _, incident := range activeIncidents {
		if incident.IsPointInRadius(lat, lng) {
			loc.IncidentIDs = append(loc.IncidentIDs, incident.ID)
		}
	}

	loc.IsCheck = loc.HasIncidents()
	if err := s.Repo.Save(loc); err != nil {
		return nil, err
	}
	if loc.IsCheck {
		s.Queue.Enqueue(loc)
	}
	return loc, nil
}
