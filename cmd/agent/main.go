package main

import (
	"log"
	"time"

	"github.com/Anton119/go-metrics-collector.git/internal/agent"
	"github.com/Anton119/go-metrics-collector.git/internal/flags"
)

func main() {
	collector := agent.NewCollector()

	flags.ParseFlags()

	myAgent := &agent.Agent{
		Collector:      collector,
		ServerAddress:  "http://localhost" + flags.FlagRunAddr,
		PollInterval:   time.Duration(flags.FlagPollInterval) * time.Second,
		ReportInterval: time.Duration(flags.FlagReportInterval) * time.Second,
	}

	log.Printf("Agent starting. PollInterval=%ds, ReportInterval=%ds, ServerAddress=%s\n",
		flags.FlagPollInterval, flags.FlagReportInterval, myAgent.ServerAddress)

	go myAgent.StartPolling()
	go myAgent.StartReporting()

	select {}
}
