package main

import (
	"net/http"

	"github.com/Anton119/go-metrics-collector.git/internal/handler"
	"github.com/Anton119/go-metrics-collector.git/internal/repository"
	"github.com/Anton119/go-metrics-collector.git/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	storage := repository.NewMemStorage()
	svc := service.NewMetricsService(storage)
	h := handler.NewMetricsHandler(svc)

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetrics)
	r.Get("/value/{type}/{name}", h.GetMetrics)
	r.Get("/", h.GetAllMetrics)

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}

}
