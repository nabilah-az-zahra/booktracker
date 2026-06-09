package services

import (
	"booktracker/backend/internal/models"
	"booktracker/backend/internal/repositories"
)

type StatsService struct {
	statsRepo *repositories.StatsRepository
}

func NewStatsService(statsRepo *repositories.StatsRepository) *StatsService {
	return &StatsService{statsRepo: statsRepo}
}

func (s *StatsService) GetStats(userID string) (*models.StatsResult, error) {
	return s.statsRepo.GetStats(userID)
}
