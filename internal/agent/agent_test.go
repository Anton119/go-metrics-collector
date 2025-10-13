package agent_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Anton119/go-metrics-collector.git/internal/agent"
	"github.com/stretchr/testify/require"
)

func TestCollector_Collect(t *testing.T) {
	c := agent.NewCollector()

	metrics := c.GetMetrics()
	require.Empty(t, metrics)

	c.Collect()
	metrics = c.GetMetrics()

	expectedGauges := map[string]bool{
		"Alloc": true, "BuckHashSys": true, "Frees": true, "GCCPUFraction": true, "GCSys": true,
		"HeapAlloc": true, "HeapIdle": true, "HeapInuse": true, "HeapObjects": true, "HeapReleased": true,
		"HeapSys": true, "LastGC": true, "Lookups": true, "MCacheInuse": true, "MCacheSys": true,
		"MSpanInuse": true, "MSpanSys": true, "Mallocs": true, "NextGC": true, "NumForcedGC": true,
		"NumGC": true, "OtherSys": true, "PauseTotalNs": true, "StackInuse": true, "StackSys": true,
		"Sys": true, "TotalAlloc": true, "RandomValue": true,
	}

	var gaugeCount, counterCount int
	for _, m := range metrics {
		if m.MType == "gauge" {
			require.True(t, expectedGauges[m.ID], "Unexpected gauge: %s", m.ID)
			gaugeCount++
		} else if m.MType == "counter" {
			require.Equal(t, "PollCount", m.ID)
			counterCount++
		}
	}

	require.Equal(t, len(expectedGauges), gaugeCount)
	require.Equal(t, 1, counterCount)
}

func TestAgent_ReportOnce(t *testing.T) {
	received := make(map[string]bool)

	// cоздаем тестовый сервер, который принимает метрики
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received[r.URL.Path] = true
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	collector := agent.NewCollector()
	agentInstance := &agent.Agent{
		Collector:      collector,
		ServerAddress:  testServer.URL,
		PollInterval:   10 * time.Millisecond,
		ReportInterval: 10 * time.Millisecond,
	}

	// cобираем метрики один раз
	collector.Collect()

	agentInstance.ReportOnce()

	time.Sleep(20 * time.Millisecond)

	require.NotEmpty(t, received, "No metrics were received by the server")
}
