package usecase

import (
	"context"

	"github.com/Soujuruya/01_SPEC/internal/domain/incident"
	"github.com/google/uuid"
)

type IncidentService struct {
	Repo incident.IncidentRepository
}

func NewIncidentService(repo incident.IncidentRepository) *IncidentService {
	return &IncidentService{Repo: repo}
}

func (s *IncidentService) CreateIncident(ctx context.Context, inc *incident.Incident) error {
	return s.Repo.Create(ctx, inc)
}

func (s *IncidentService) GetIncident(ctx context.Context, id uuid.UUID) (*incident.Incident, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *IncidentService) UpdateIncident(ctx context.Context, inc *incident.Incident) error {
	return s.Repo.Update(ctx, inc)
}

func (s *IncidentService) DeactivateIncident(ctx context.Context, id uuid.UUID) error {
	return s.Repo.Deactivate(ctx, id)
}

func (s *IncidentService) ListIncidents(ctx context.Context, offset, limit int) ([]*incident.Incident, int, error) {
	return s.Repo.ListWithTotal(ctx, offset, limit)
}

func (s *IncidentService) GetActiveIncidents(ctx context.Context) ([]*incident.Incident, error) {
	return s.Repo.GetActiveIncidents(ctx)
}

func (s *IncidentService) CountActiveIncidents(ctx context.Context) (int, error) {
	return s.Repo.CountActiveIncidents(ctx)
}
