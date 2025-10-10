package handler

import (
	"fmt"
	"net/http"
	"strings"

	models "github.com/Anton119/go-metrics-collector.git/internal/model"
	"github.com/Anton119/go-metrics-collector.git/internal/service"
)

type MetricsHandler struct {
	svc *service.MetricsService
}

func NewMetricsHandler(svc *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		svc: svc,
	}
}

func (h *MetricsHandler) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "неверный метод запроса", http.StatusBadRequest)
		return
	}

	err := h.svc.UpdateMetrics(r.URL.Path)
	if err != nil {

		switch err.Error() {
		case "имя метрики не указано", "неверный формат пути":
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("ok"))

}

func (h *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "неверный метод запроса", http.StatusBadRequest)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 {
		http.Error(w, "неверный формат пути", http.StatusBadRequest)
		return
	}

	mType := parts[1]
	id := parts[2]

	metrics, err := h.svc.GetMetrics(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if metrics.MType != mType {
		http.Error(w, "тип метрики не совпадает", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if mType == models.Gauge {
		w.Write([]byte(fmt.Sprintf("%g", *metrics.Value)))
	} else if mType == models.Counter {
		w.Write([]byte(fmt.Sprintf("%d", *metrics.Delta)))
	} else {
		w.Write([]byte("значение не установлено"))
	}
}
