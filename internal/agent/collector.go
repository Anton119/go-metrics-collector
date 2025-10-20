package agent

import (
	"log"
	"math/rand"
	"runtime"
	"sync"

	models "github.com/Anton119/go-metrics-collector.git/internal/model"
)

type Collector struct {
	mu       *sync.Mutex
	gauges   map[string]float64
	counters map[string]int64
}

func NewCollector() *Collector {
	return &Collector{
		mu:       &sync.Mutex{},
		gauges:   make(map[string]float64),
		counters: map[string]int64{},
	}
}

func (c *Collector) Collect() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	randomValue := rand.Float64()

	newGauges := map[string]float64{
		"Alloc":         float64(m.Alloc),
		"BuckHashSys":   float64(m.BuckHashSys),
		"Frees":         float64(m.Frees),
		"GCCPUFraction": m.GCCPUFraction,
		"GCSys":         float64(m.GCSys),
		"HeapAlloc":     float64(m.HeapAlloc),
		"HeapIdle":      float64(m.HeapIdle),
		"HeapInuse":     float64(m.HeapInuse),
		"HeapObjects":   float64(m.HeapObjects),
		"HeapReleased":  float64(m.HeapReleased),
		"HeapSys":       float64(m.HeapSys),
		"LastGC":        float64(m.LastGC),
		"Lookups":       float64(m.Lookups),
		"MCacheInuse":   float64(m.MCacheInuse),
		"MCacheSys":     float64(m.MCacheSys),
		"MSpanInuse":    float64(m.MSpanInuse),
		"MSpanSys":      float64(m.MSpanSys),
		"Mallocs":       float64(m.Mallocs),
		"NextGC":        float64(m.NextGC),
		"NumForcedGC":   float64(m.NumForcedGC),
		"NumGC":         float64(m.NumGC),
		"OtherSys":      float64(m.OtherSys),
		"PauseTotalNs":  float64(m.PauseTotalNs),
		"StackInuse":    float64(m.StackInuse),
		"StackSys":      float64(m.StackSys),
		"Sys":           float64(m.Sys),
		"TotalAlloc":    float64(m.TotalAlloc),
		"RandomValue":   randomValue,
	}

	c.mu.Lock()
	for k, v := range newGauges {
		c.gauges[k] = v
	}
	c.counters["PollCount"]++
	c.mu.Unlock()

	log.Printf("[Collector] Собрано метрик: gauges=%d, counters=%d\n", len(c.gauges), len(c.counters))
}

func (c *Collector) GetMetrics() []models.Metrics {
	c.mu.Lock()
	defer c.mu.Unlock()

	var result []models.Metrics

	for k, v := range c.gauges {
		val := &v
		result = append(result, models.Metrics{
			ID:    k,
			MType: models.Gauge,
			Value: val,
		})
	}

	for k, v := range c.counters {
		d := &v
		result = append(result, models.Metrics{
			ID:    k,
			MType: models.Counter,
			Delta: d,
		})
	}

	return result
}
