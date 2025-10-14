package main

import (
	"time"

	"github.com/Anton119/go-metrics-collector.git/internal/agent"
)

func main() {
	collector := agent.NewCollector()

	myAgent := &agent.Agent{
		Collector:      collector,
		ServerAddress:  "http://localhost:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
	}

	go myAgent.StartPolling()
	go myAgent.StartReporting()

	select {}
}
