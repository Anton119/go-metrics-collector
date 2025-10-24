package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Anton119/go-metrics-collector.git/internal/repository"
	"github.com/Anton119/go-metrics-collector.git/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

// setupRouter создаёт готовый Router с хендлерами и хранилищем
func setupRouter() (*chi.Mux, *MetricsHandler) {
	storage := repository.NewMemStorage()
	svc := service.NewMetricsService(storage)
	h := NewMetricsHandler(svc)

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetrics)
	r.Get("/value/{type}/{name}", h.GetMetrics)
	r.Get("/", h.GetAllMetrics)

	return r, h
}

func TestUpdateMetrics(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/TestMetric/42", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetMetrics(t *testing.T) {
	r, _ := setupRouter()

	reqUpdate := httptest.NewRequest(http.MethodPost, "/update/gauge/TestMetric/42", nil)
	wUpdate := httptest.NewRecorder()
	r.ServeHTTP(wUpdate, reqUpdate)

	reqGet := httptest.NewRequest(http.MethodGet, "/value/gauge/TestMetric", nil)
	wGet := httptest.NewRecorder()
	r.ServeHTTP(wGet, reqGet)

	respGet := wGet.Result()
	defer respGet.Body.Close()
	require.Equal(t, http.StatusOK, respGet.StatusCode)
}

func TestGetAllMetrics(t *testing.T) {
	r, _ := setupRouter()

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/update/gauge/TestMetric1/10", nil))
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/update/counter/TestMetric2/5", nil))

	reqAll := httptest.NewRequest(http.MethodGet, "/", nil)
	wAll := httptest.NewRecorder()
	r.ServeHTTP(wAll, reqAll)

	respAll := wAll.Result()
	defer respAll.Body.Close()
	require.Equal(t, http.StatusOK, respAll.StatusCode)
}
