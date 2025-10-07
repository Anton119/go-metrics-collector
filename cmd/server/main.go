package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type Metric struct {
	Type    MetricType
	Gauge   *float64
	Counter *int64
}

type MemStorage struct {
	metric map[string]Metric
}

func UpdateMetrics(w http.ResponseWriter, r *http.Request, s *MemStorage) {
	if r.Method != http.MethodPost {
		http.Error(w, "неверный метод запроса", http.StatusBadRequest)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) != 4 {
		http.Error(w, "неверный формат пути", http.StatusNotFound)
		return
	}

	metricType := parts[1]
	metricName := parts[2]
	metricValue := parts[3]

	if metricName == "" {
		http.Error(w, "имя метрики не указано", http.StatusNotFound)
		return
	}

	switch metricType {
	case string(Gauge):
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "некорректное значение для gauge", http.StatusBadRequest)
			return
		}

		s.metric[metricName] = Metric{
			Type:    Gauge,
			Gauge:   &val,
			Counter: nil,
		}

	case string(Counter):
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "некорректное значение для counter", http.StatusBadRequest)
			return
		}

		s.metric[metricName] = Metric{
			Type:    Counter,
			Gauge:   nil,
			Counter: &val,
		}

	default:
		http.Error(w, "неизвестный тип метрики", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("ok"))

}

func main() {

	storage := &MemStorage{
		metric: make(map[string]Metric),
	}

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, func(w http.ResponseWriter, r *http.Request) {
		UpdateMetrics(w, r, storage)
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
