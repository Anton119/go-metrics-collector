package main

import (
	"log"
	"strings"
	"time"

	"github.com/Anton119/go-metrics-collector.git/internal/agent"
	"github.com/Anton119/go-metrics-collector.git/internal/flags"
)

func main() {
	flags.ParseFlags()

	collector := agent.NewCollector()

	serverAddress := flags.FlagRunAddr
	if !strings.HasPrefix(serverAddress, "http://") && !strings.HasPrefix(serverAddress, "https://") {
		serverAddress = "http://" + serverAddress
	}

	myAgent := &agent.Agent{
		Collector:      collector,
		ServerAddress:  serverAddress,
		PollInterval:   time.Duration(flags.FlagPollInterval) * time.Second,
		ReportInterval: time.Duration(flags.FlagReportInterval) * time.Second,
	}

	log.Printf("Agent starting. PollInterval=%ds, ReportInterval=%ds, ServerAddress=%s\n",
		flags.FlagPollInterval, flags.FlagReportInterval, myAgent.ServerAddress)

	go myAgent.StartPolling()
	go myAgent.StartReporting()

	select {} // блокируем main
}
