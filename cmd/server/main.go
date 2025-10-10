package main

import (
	"net/http"

	"github.com/Anton119/go-metrics-collector.git/internal/handler"
	"github.com/Anton119/go-metrics-collector.git/internal/repository"
	"github.com/Anton119/go-metrics-collector.git/internal/service"
)

func main() {
	storage := repository.NewMemStorage()
	svc := service.NewMetricsService(storage)
	h := handler.NewMetricsHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", h.UpdateMetrics)
	mux.HandleFunc("/value/", h.GetMetrics)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}

}
