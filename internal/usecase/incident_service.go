package usecase

import (
	"context"

	"github.com/Soujuruya/01_SPEC/internal/domain/incident"
	"github.com/Soujuruya/01_SPEC/internal/pkg/logger"
	"github.com/google/uuid"
)

type IncidentService struct {
	Repo  incident.IncidentRepository
	Cache incident.IncidentCache
	lg    *logger.Logger
}

func NewIncidentService(repo incident.IncidentRepository, cache incident.IncidentCache, lg *logger.Logger) *IncidentService {
	return &IncidentService{
		Repo:  repo,
		Cache: cache,
		lg:    lg,
	}
}

func (s *IncidentService) CreateIncident(ctx context.Context, inc *incident.Incident) error {
	if err := s.Repo.Create(ctx, inc); err != nil {
		s.lg.Error("CreateIncident failed", "incident_id", inc.ID, "error", err)
		return err
	}

	_ = s.Cache.InvalidateActive(ctx)
	s.lg.Info("Created incident", "incident_id", inc.ID, "title", inc.Title)
	return nil
}

func (s *IncidentService) GetIncident(ctx context.Context, id uuid.UUID) (*incident.Incident, error) {
	inc, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		s.lg.Error("GetIncident failed", "incident_id", id, "error", err)
		return nil, err
	}
	s.lg.Debug("GetIncident success", "incident_id", id)
	return inc, nil
}

func (s *IncidentService) UpdateIncident(ctx context.Context, inc *incident.Incident) error {
	if err := s.Repo.Update(ctx, inc); err != nil {
		s.lg.Error("UpdateIncident failed", "incident_id", inc.ID, "error", err)
		return err
	}

	_ = s.Cache.InvalidateActive(ctx)
	s.lg.Info("Updated incident", "incident_id", inc.ID)
	return nil
}

func (s *IncidentService) DeactivateIncident(ctx context.Context, id uuid.UUID) error {
	if err := s.Repo.Deactivate(ctx, id); err != nil {
		s.lg.Error("DeactivateIncident failed", "incident_id", id, "error", err)
		return err
	}

	_ = s.Cache.InvalidateActive(ctx)
	s.lg.Info("Deactivated incident", "incident_id", id)
	return nil
}

func (s *IncidentService) ListIncidents(ctx context.Context, offset, limit int) ([]*incident.Incident, int, error) {
	incs, total, err := s.Repo.ListWithTotal(ctx, offset, limit)
	if err != nil {
		s.lg.Error("ListIncidents failed", "offset", offset, "limit", limit, "error", err)
		return nil, 0, err
	}
	s.lg.Debug("ListIncidents success", "offset", offset, "limit", limit, "returned", len(incs), "total", total)
	return incs, total, nil
}

func (s *IncidentService) GetActiveIncidents(ctx context.Context) ([]*incident.Incident, error) {
	incs, err := s.Cache.GetActive(ctx)
	if err == nil {
		s.lg.Debug("GetActiveIncidents cache hit", "count", len(incs))
		return incs, nil
	}

	incs, err = s.Repo.GetActiveIncidents(ctx)
	if err != nil {
		s.lg.Error("GetActiveIncidents failed", "error", err)
		return nil, err
	}

	_ = s.Cache.SetActive(ctx, incs)
	s.lg.Debug("GetActiveIncidents cache set", "count", len(incs))
	return incs, nil
}

func (s *IncidentService) CountActiveIncidents(ctx context.Context) (int, error) {
	count, err := s.Repo.CountActiveIncidents(ctx)
	if err != nil {
		s.lg.Error("CountActiveIncidents failed", "error", err)
		return 0, err
	}
	s.lg.Debug("CountActiveIncidents success", "count", count)
	return count, nil
}
