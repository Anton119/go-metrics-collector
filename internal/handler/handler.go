package handler

import (
	"fmt"
	"net/http"

	models "github.com/Anton119/go-metrics-collector.git/internal/model"
	"github.com/Anton119/go-metrics-collector.git/internal/service"
	"github.com/go-chi/chi/v5"
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
	fmt.Printf("[Server] UpdateMetrics called: %s %s\n", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "неверный метод запроса", http.StatusBadRequest)
		return
	}

	err := h.svc.UpdateMetrics(r.URL.Path)
	if err != nil {
		fmt.Printf("[Server] Error updating metric: %s, %v\n", r.URL.Path, err)

		switch err.Error() {
		case "имя метрики не указано", "неверный формат пути":
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Printf("[Server] Metric updated successfully: %s\n", r.URL.Path)
	w.Write([]byte("ok"))

}

func (h *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[Server] GetMetrics called: %s %s\n", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "неверный метод запроса", http.StatusBadRequest)
		return
	}

	mType := chi.URLParam(r, "type")
	id := chi.URLParam(r, "name")

	metrics, err := h.svc.GetMetrics(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if metrics.MType != mType {
		http.Error(w, "тип метрики не совпадает", http.StatusBadRequest)
		return
	}
	fmt.Printf("[Server] Returning metric: %s = ", id)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if mType == models.Gauge {
		w.Write([]byte(fmt.Sprintf("%g", *metrics.Value)))
	} else if mType == models.Counter {
		w.Write([]byte(fmt.Sprintf("%d", *metrics.Delta)))
	} else {
		w.Write([]byte("значение не установлено"))
	}
}

func (h *MetricsHandler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "неверный метод запроса", http.StatusBadRequest)
		return
	}

	allMetrics, err := h.svc.GetAllMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//пишет строку в w и добавляет перенос строки
	fmt.Fprintln(w, "<html><body><ul>")

	for _, m := range allMetrics {
		var value string
		if m.MType == models.Gauge {
			value = fmt.Sprintf("%g", *m.Value)
		} else {
			value = fmt.Sprintf("%d", *m.Delta)
		}
		fmt.Fprintf(w, "<li>%s (%s) = %s</li>", m.ID, m.MType, value)
	}

	fmt.Fprintln(w, "</ul></body></html>")
}
