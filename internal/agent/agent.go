package agent

import (
	"fmt"
	"net/http"
	"time"

	models "github.com/Anton119/go-metrics-collector.git/internal/model"
)

type Agent struct {
	Collector      *Collector
	ServerAddress  string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func (a *Agent) StartPolling() {
	ticker := time.NewTicker(a.PollInterval)
	defer ticker.Stop()

	for range ticker.C {
		a.Collector.Collect()
	}
}

func (a *Agent) StartReporting() {
	ticker := time.NewTicker(a.ReportInterval)
	defer ticker.Stop()

	for range ticker.C {
		metrics := a.Collector.GetMetrics()
		fmt.Printf("[Agent] Отправка %d метрик на сервер %s\n", len(metrics), a.ServerAddress)

		for _, m := range metrics {
			var value string
			if m.MType == models.Gauge {
				value = fmt.Sprintf("%f", *m.Value)
			} else {
				value = fmt.Sprintf("%d", *m.Delta)
			}
			url := fmt.Sprintf("%s/update/%s/%s/%s", a.ServerAddress, m.MType, m.ID, value)
			_, err := http.Post(url, "text/plain", nil)
			if err != nil {
				fmt.Printf("Error sending metric %s: %v\n", m.ID, err)
			}
		}
	}

}

// для юнит теста
func (a *Agent) ReportOnce() {
	metrics := a.Collector.GetMetrics()
	for _, m := range metrics {
		var value string
		if m.MType == models.Gauge {
			value = fmt.Sprintf("%f", *m.Value)
		} else {
			value = fmt.Sprintf("%d", *m.Delta)
		}

		http.Post(
			fmt.Sprintf("%s/update/%s/%s/%s", a.ServerAddress, m.MType, m.ID, value),
			"text/plain",
			nil,
		)
	}
}
