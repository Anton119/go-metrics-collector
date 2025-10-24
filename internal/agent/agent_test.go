package agent_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
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
		switch m.MType {
		case "gauge":
			require.True(t, expectedGauges[m.ID], "Unexpected gauge: %s", m.ID)
			gaugeCount++
		case "counter":
			require.Equal(t, "PollCount", m.ID)
			counterCount++
		default:
			t.Errorf("Unknown metric type: %s", m.MType)
		}
	}

	require.Equal(t, len(expectedGauges), gaugeCount)
	require.Equal(t, 1, counterCount)
}

func TestAgent_ReportOnce(t *testing.T) {
	received := make(map[string]bool)

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

	collector.Collect()

	agentInstance.ReportOnce()

	require.NotEmpty(t, received, "No metrics were received by the server")

	found := false
	for path := range received {
		if pathContainsMetric(path, "RandomValue") || pathContainsMetric(path, "PollCount") {
			found = true
			break
		}
	}
	require.True(t, found, "Expected metrics not received")
}

func pathContainsMetric(path, metricID string) bool {
	return strings.Contains(path, metricID)
}
