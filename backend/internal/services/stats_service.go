package services

import (
	"booktracker/backend/internal/models"
	redisclient "booktracker/backend/internal/redis"
	"booktracker/backend/internal/repositories"
	"context"
	"encoding/json"
	"time"
)

type StatsService struct {
	statsRepo *repositories.StatsRepository
}

func NewStatsService(statsRepo *repositories.StatsRepository) *StatsService {
	return &StatsService{statsRepo: statsRepo}
}

func (s *StatsService) GetStats(ctx context.Context, userID string) (*models.StatsResult, error) {
    cacheKey := "stats:" + userID
    
    cached, err := redisclient.Client.Get(ctx, cacheKey).Result()
    if err == nil {
        var stats models.StatsResult
        if json.Unmarshal([]byte(cached), &stats) == nil {
            return &stats, nil
        }
    }
    
    stats, err := s.statsRepo.GetStats(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if data, err := json.Marshal(stats); err == nil {
        redisclient.Client.Set(ctx, cacheKey, data, 60*time.Second)
    }
    
    return stats, nil
}

func (s *StatsService) GetReadingHistory(ctx context.Context, userID string, days int) ([]models.DailyReading, error) {
	return s.statsRepo.GetReadingHistory(ctx, userID, days)
}