package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Anton119/go-metrics-collector.git/internal/repository"
	"github.com/Anton119/go-metrics-collector.git/internal/service"
	"github.com/stretchr/testify/require"
)

func TestUpdateAndGetMetrics(t *testing.T) {
	storage := repository.NewMemStorage()
	svc := service.NewMetricsService(storage)
	h := NewMetricsHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/TestMetric/42", nil)
	w := httptest.NewRecorder()
	h.UpdateMetrics(w, req)
	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	reqGet := httptest.NewRequest(http.MethodGet, "/value/gauge/TestMetric", nil)
	wGet := httptest.NewRecorder()
	h.GetMetrics(w, reqGet)
	respGet := wGet.Result()
	require.Equal(t, http.StatusOK, respGet.StatusCode)

}
