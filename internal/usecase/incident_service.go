package usecase

import (
	"context"

	"github.com/Soujuruya/01_SPEC/internal/domain/incident"
	"github.com/google/uuid"
)

type IncidentService struct {
	Repo  incident.IncidentRepository
	Cache incident.IncidentCache
}

func NewIncidentService(repo incident.IncidentRepository, cache incident.IncidentCache) *IncidentService {
	return &IncidentService{
		Repo:  repo,
		Cache: cache,
	}
}

func (s *IncidentService) CreateIncident(ctx context.Context, inc *incident.Incident) error {
	if err := s.Repo.Create(ctx, inc); err != nil {
		return err
	}

	_ = s.Cache.InvalidateActive(ctx)
	return nil
}

func (s *IncidentService) GetIncident(ctx context.Context, id uuid.UUID) (*incident.Incident, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *IncidentService) UpdateIncident(ctx context.Context, inc *incident.Incident) error {
	if err := s.Repo.Update(ctx, inc); err != nil {
		return err
	}

	_ = s.Cache.InvalidateActive(ctx)
	return nil
}

func (s *IncidentService) DeactivateIncident(ctx context.Context, id uuid.UUID) error {
	if err := s.Repo.Deactivate(ctx, id); err != nil {
		return err
	}

	_ = s.Cache.InvalidateActive(ctx)
	return nil
}

func (s *IncidentService) ListIncidents(ctx context.Context, offset, limit int) ([]*incident.Incident, int, error) {
	return s.Repo.ListWithTotal(ctx, offset, limit)
}

func (s *IncidentService) GetActiveIncidents(ctx context.Context) ([]*incident.Incident, error) {
	incs, err := s.Cache.GetActive(ctx)
	if err == nil {
		return incs, nil
	}

	incs, err = s.Repo.GetActiveIncidents(ctx)
	if err != nil {
		return nil, err
	}

	_ = s.Cache.SetActive(ctx, incs)

	return incs, nil
}

func (s *IncidentService) CountActiveIncidents(ctx context.Context) (int, error) {
	return s.Repo.CountActiveIncidents(ctx)
}
