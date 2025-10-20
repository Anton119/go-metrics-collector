package handler

import (
	"github.com/Anton119/go-metrics-collector.git/internal/service"
	"github.com/go-chi/chi/v5"
)

func NewRouter(svc *service.MetricsService) *chi.Mux {
	h := NewMetricsHandler(svc)
	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", h.UpdateMetrics)
	r.Get("/value/{type}/{name}", h.GetMetrics)
	r.Get("/", h.GetAllMetrics)

	return r
}
