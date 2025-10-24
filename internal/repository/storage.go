package repository

import models "github.com/Anton119/go-metrics-collector.git/internal/model"

type MemStorage struct {
	metrics map[string]models.Metrics
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]models.Metrics),
	}
}

func (s *MemStorage) SetMetrics(metric models.Metrics) {
	s.metrics[metric.ID] = metric
}

func (s *MemStorage) GetMetrics(id string) (models.Metrics, bool) {
	m, ok := s.metrics[id]
	return m, ok
}

func (s *MemStorage) GetAllMetrics() []models.Metrics {
	metrics := make([]models.Metrics, 0, len(s.metrics))
	for _, m := range s.metrics {
		metrics = append(metrics, m)
	}

	return metrics
}
