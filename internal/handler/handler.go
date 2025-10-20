package handler

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"

	models "github.com/Anton119/go-metrics-collector.git/internal/model"
	"github.com/Anton119/go-metrics-collector.git/internal/service"
	"github.com/go-chi/chi/v5"
)

var (
	ErrEmptyMetricName = errors.New("имя метрики не указано")
	ErrInvalidValue    = errors.New("неверный формат значения")
	ErrInvalidType     = errors.New("неверный тип метрики")
	ErrMetricNotFound  = errors.New("метрика не найдена")
	ErrNoMetrics       = errors.New("метрики не найдены")
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
	log.Printf("[Server] UpdateMetrics called: %s %s\n", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "неверный метод запроса", http.StatusBadRequest)
		return
	}

	mType := chi.URLParam(r, "type")
	id := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	err := h.svc.UpdateMetrics(mType, id, value)
	if err != nil {
		log.Printf("[Server] Error updating metric: %s/%s = %s, %v\n", mType, id, value, err)

		switch {
		case errors.Is(err, ErrEmptyMetricName):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrInvalidValue):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrInvalidType):
			http.Error(w, err.Error(), http.StatusNotImplemented)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	log.Printf("[Server] Metric updated successfully: %s/%s = %s\n", mType, id, value)
	w.Write([]byte("ok"))

}

func (h *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Server] GetMetrics called: %s %s\n", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "неверный метод запроса", http.StatusBadRequest)
		return
	}

	mType := chi.URLParam(r, "type")
	id := chi.URLParam(r, "name")

	metrics, err := h.svc.GetMetrics(id)
	if err != nil {
		if errors.Is(err, ErrMetricNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if metrics.MType != mType {
		http.Error(w, "тип метрики не совпадает", http.StatusBadRequest)
		return
	}
	log.Printf("[Server] Returning metric: %s = ", id)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if mType == models.Gauge {
		w.Write([]byte(strconv.FormatFloat(*metrics.Value, 'f', -1, 64)))
	} else if mType == models.Counter {
		w.Write([]byte(strconv.FormatInt(*metrics.Delta, 10)))
	} else {
		w.Write([]byte("значение не установлено"))
	}
}

type MetricView struct {
	ID    string
	MType string
	Value string
}

var metricsTemplate = template.Must(template.New("metrics").Parse(`
<html>
  <body>
    <ul>
      {{range .}}
        <li>{{.ID}} ({{.MType}}) = {{.Value}}</li>
      {{end}}
    </ul>
  </body>
</html>
`))

func (h *MetricsHandler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "неверный метод запроса", http.StatusBadRequest)
		return
	}

	allMetrics, err := h.svc.GetAllMetrics()
	if err != nil {
		if errors.Is(err, ErrNoMetrics) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var viewData []MetricView
	for _, m := range allMetrics {
		var val string
		if m.MType == models.Gauge {
			val = strconv.FormatFloat(*m.Value, 'f', -1, 64)
		} else {
			val = strconv.FormatInt(*m.Delta, 10)
		}
		viewData = append(viewData, MetricView{
			ID:    m.ID,
			MType: m.MType,
			Value: val,
		})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := metricsTemplate.Execute(w, viewData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
