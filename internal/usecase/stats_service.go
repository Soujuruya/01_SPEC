package usecase

import (
	"context"
	"time"

	"github.com/Soujuruya/01_SPEC/internal/domain/location"
	"github.com/Soujuruya/01_SPEC/internal/pkg/logger"
)

type StatsService struct {
	Repo location.LocationRepository
	Lg   *logger.Logger
}

func NewStatsService(repo location.LocationRepository, lg *logger.Logger) *StatsService {
	return &StatsService{
		Repo: repo,
		Lg:   lg,
	}
}

func (s *StatsService) GetUserCount(ctx context.Context, minutes int) (int, error) {
	since := time.Now().Add(-time.Duration(minutes) * time.Minute)
	count, err := s.Repo.CountUniqueUsers(ctx, since)
	if err != nil {
		s.Lg.Error("StatsService.GetUserCount: failed to count unique users", "error", err, "minutes", minutes)
		return 0, err
	}
	s.Lg.Debug("StatsService.GetUserCount: unique user count fetched", "count", count, "minutes", minutes)
	return count, nil
}
