package service

import (
	"errors"
	"strconv"

	models "github.com/Anton119/go-metrics-collector.git/internal/model"
	"github.com/Anton119/go-metrics-collector.git/internal/repository"
)

var (
	ErrEmptyMetricName = errors.New("имя метрики не указано")
	ErrInvalidValue    = errors.New("неверное значение")
	ErrInvalidType     = errors.New("неизвестный тип метрики")
	ErrMetricNotFound  = errors.New("метрика не найдена")
	ErrNoMetrics       = errors.New("метрики не найдены")
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
		return ErrEmptyMetricName
	}

	metric := models.Metrics{
		ID:    id,
		MType: mType,
	}

	switch mType {
	case models.Gauge:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return ErrInvalidValue
		}
		metric.Value = &f

	case models.Counter:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return ErrInvalidValue
		}

		if existing, ok := s.storage.GetMetrics(id); ok && existing.Delta != nil {
			i += *existing.Delta
		}
		metric.Delta = &i

	default:
		return ErrInvalidType
	}

	s.storage.SetMetrics(metric)
	return nil
}

func (s *MetricsService) GetMetrics(id string) (models.Metrics, error) {
	m, ok := s.storage.GetMetrics(id)
	if !ok {
		return models.Metrics{}, ErrMetricNotFound
	}
	return m, nil
}

func (s *MetricsService) GetAllMetrics() ([]models.Metrics, error) {
	metrics := s.storage.GetAllMetrics()
	if len(metrics) == 0 {
		return nil, ErrNoMetrics
	}

	return metrics, nil
}
