package usecase

import (
	"context"
	"time"

	"github.com/Soujuruya/01_SPEC/internal/domain/incident"
	"github.com/Soujuruya/01_SPEC/internal/domain/location"
	"github.com/google/uuid"
)

type LocationService struct {
	Repo         location.LocationRepository
	IncidentRepo incident.IncidentRepository
	Queue        location.WebhookQueue
}

func NewLocationService(repo location.LocationRepository, incidentRepo incident.IncidentRepository, queue location.WebhookQueue) *LocationService {
	return &LocationService{
		Repo:         repo,
		IncidentRepo: incidentRepo,
		Queue:        queue,
	}
}

// CheckLocation проверяет координаты пользователя
func (s *LocationService) CheckLocation(ctx context.Context, userID uuid.UUID, lat, lng float64) (*location.Location, error) {
	activeIncidents, _ := s.IncidentRepo.GetActiveIncidents(ctx)

	loc := &location.Location{
		ID:          uuid.New(),
		UserID:      userID,
		Lat:         lat,
		Lng:         lng,
		Timestamp:   time.Now(),
		IncidentIDs: []uuid.UUID{},
	}

	for _, inc := range activeIncidents {
		if inc.IsPointInRadius(lat, lng) {
			loc.IncidentIDs = append(loc.IncidentIDs, inc.ID)
		}
	}

	loc.IsCheck = len(loc.IncidentIDs) > 0

	if err := s.Repo.Save(ctx, loc); err != nil {
		return nil, err
	}

	if loc.IsCheck {
		s.Queue.Enqueue(loc)
	}

	return loc, nil
}
