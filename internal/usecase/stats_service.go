package usecase

import (
	"context"
	"time"

	"github.com/Soujuruya/01_SPEC/internal/domain/location"
)

type StatsService struct {
	Repo location.LocationRepository
}

func NewStatsService(repo location.LocationRepository) *StatsService {
	return &StatsService{Repo: repo}
}

func (s *StatsService) GetUserCount(ctx context.Context, minutes int) (int, error) {
	since := time.Now().Add(-time.Duration(minutes) * time.Minute)
	return s.Repo.CountUniqueUsers(ctx, since)
}
