package service

import (
	"errors"
	"fmt"
	"strconv"

	models "github.com/Anton119/go-metrics-collector.git/internal/model"
	"github.com/Anton119/go-metrics-collector.git/internal/repository"
)

type MetricsService struct {
	storage *repository.MemStorage
}

func NewMetricsService(storage *repository.MemStorage) *MetricsService {
	return &MetricsService{
		storage: storage,
	}
}

func (s *MetricsService) UpdateMetrics(mType, id, val string) error {
	if id == "" {
		return errors.New("имя метрики не указано")
	}

	metric := models.Metrics{
		ID:    id,
		MType: mType,
	}

	switch mType {
	case models.Gauge:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return errors.New("неверное значение для gauge")
		}
		metric.Value = &f

	case models.Counter:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return errors.New("неверное значение для counter")
		}
		metric.Delta = &i

	default:
		return errors.New("неизвестный тип метрики")
	}

	s.storage.SetMetrics(metric)
	return nil
}

func (s *MetricsService) GetMetrics(id string) (models.Metrics, error) {
	m, ok := s.storage.GetMetrics(id)
	if !ok {
		return models.Metrics{}, fmt.Errorf("метрика %s не найдена", id)
	}
	return m, nil
}

func (s *MetricsService) GetAllMetrics() ([]models.Metrics, error) {
	metrics := s.storage.GetAllMetrics()
	if len(metrics) == 0 {
		return []models.Metrics{}, fmt.Errorf("метрики не найдены")
	}

	return metrics, nil
}
