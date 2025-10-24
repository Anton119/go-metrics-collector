package main

import (
	"log"
	"net/http"

	"github.com/Anton119/go-metrics-collector.git/internal/flags"
	"github.com/Anton119/go-metrics-collector.git/internal/handler"
	"github.com/Anton119/go-metrics-collector.git/internal/repository"
	"github.com/Anton119/go-metrics-collector.git/internal/service"
)

func main() {

	flags.ParseFlags()

	storage := repository.NewMemStorage()
	svc := service.NewMetricsService(storage)
	r := handler.NewRouter(svc)

	log.Println("Server is starting on", flags.FlagRunAddr)
	if err := http.ListenAndServe(flags.FlagRunAddr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
