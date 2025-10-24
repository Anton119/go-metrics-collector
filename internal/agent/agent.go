package agent

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	models "github.com/Anton119/go-metrics-collector.git/internal/model"
)

type Agent struct {
	Collector      *Collector
	ServerAddress  string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func (a *Agent) StartPolling(ctx context.Context) {
	ticker := time.NewTicker(a.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.Collector.Collect()
		case <-ctx.Done():
			return
		}
	}
}

func (a *Agent) StartReporting(ctx context.Context) {
	ticker := time.NewTicker(a.ReportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.reportMetricsOnce()
		case <-ctx.Done():
			return
		}
	}
}

// для юнит тестов
func (a *Agent) ReportOnce() {
	a.reportMetricsOnce()
}

func (a *Agent) reportMetricsOnce() {
	metrics := a.Collector.GetMetrics()

	for _, m := range metrics {
		var value string
		switch m.MType {
		case models.Gauge:
			value = strconv.FormatFloat(*m.Value, 'f', 6, 64)
		case models.Counter:
			value = strconv.FormatInt(*m.Delta, 10)
		default:
			log.Printf("Unknown metric type: %s\n", m.MType)
			continue
		}

		fullURL, err := url.JoinPath(a.ServerAddress, "update", string(m.MType), m.ID, value)
		if err != nil {
			log.Printf("Error forming URL for metric %s: %v\n", m.ID, err)
			continue
		}

		resp, err := http.Post(fullURL, "text/plain", nil)
		if err != nil {
			log.Printf("Error sending metric %s: %v\n", m.ID, err)
			continue
		}

		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
