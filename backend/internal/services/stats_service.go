package services

import (
	"booktracker/backend/internal/models"
	"booktracker/backend/internal/repositories"
	"context"
	"fmt"
)

type StatsService struct {
	statsRepo *repositories.StatsRepository
}

func NewStatsService(statsRepo *repositories.StatsRepository) *StatsService {
	return &StatsService{statsRepo: statsRepo}
}

func (s *StatsService) GetStats(ctx context.Context, userID string) (*models.StatsResult, error) {
	return s.statsRepo.GetStats(ctx, userID)
}

func (s *StatsService) GetReadingHistory(ctx context.Context, userID string, days int) ([]models.DailyReading, error) {
	if days != 7 && days != 30 && days != 90 {
		return nil, fmt.Errorf("days must be 7, 30 or 90")
	}
	return s.statsRepo.GetReadingHistory(ctx, userID, days)
}