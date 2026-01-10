package usecase

import (
	"context"
	"time"

	"github.com/Soujuruya/01_SPEC/internal/domain/incident"
	"github.com/Soujuruya/01_SPEC/internal/domain/location"
	"github.com/Soujuruya/01_SPEC/internal/pkg/logger"
	"github.com/google/uuid"
)

type LocationService struct {
	Repo         location.LocationRepository
	IncidentRepo incident.IncidentRepository
	Queue        location.WebhookQueue
	Lg           *logger.Logger
}

func NewLocationService(
	repo location.LocationRepository,
	incidentRepo incident.IncidentRepository,
	queue location.WebhookQueue,
	lg *logger.Logger,
) *LocationService {
	return &LocationService{
		Repo:         repo,
		IncidentRepo: incidentRepo,
		Queue:        queue,
		Lg:           lg,
	}
}

// CheckLocation проверяет координаты пользователя
func (s *LocationService) CheckLocation(ctx context.Context, userID uuid.UUID, lat, lng float64) (*location.Location, error) {
	activeIncidents, err := s.IncidentRepo.GetActiveIncidents(ctx)
	if err != nil {
		s.Lg.Error("LocationService.CheckLocation: failed to get active incidents", "error", err, "user_id", userID)
		return nil, err
	}

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
		s.Lg.Error("LocationService.CheckLocation: failed to save location", "error", err, "user_id", userID, "location_id", loc.ID)
		return nil, err
	}

	s.Lg.Debug("LocationService.CheckLocation: location saved", "user_id", userID, "location_id", loc.ID, "incidents_found", len(loc.IncidentIDs))

	if loc.IsCheck {
		if err := s.Queue.Enqueue(ctx, loc); err != nil {
			s.Lg.Error("LocationService.CheckLocation: failed to enqueue webhook", "error", err, "user_id", userID, "location_id", loc.ID)
			return nil, err
		}
		s.Lg.Debug("LocationService.CheckLocation: webhook enqueued", "user_id", userID, "location_id", loc.ID)
	}

	return loc, nil
}
