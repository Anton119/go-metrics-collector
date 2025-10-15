package agent

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"

	models "github.com/Anton119/go-metrics-collector.git/internal/model"
)

type Collector struct {
	mu       sync.Mutex
	gauges   map[string]float64
	counters map[string]int64
}

func NewCollector() *Collector {
	return &Collector{
		gauges:   make(map[string]float64),
		counters: map[string]int64{},
	}
}

func (c *Collector) Collect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Метрики типа Gauge (float64)
	c.gauges["Alloc"] = float64(m.Alloc)
	c.gauges["BuckHashSys"] = float64(m.BuckHashSys)
	c.gauges["Frees"] = float64(m.Frees)
	c.gauges["GCCPUFraction"] = m.GCCPUFraction
	c.gauges["GCSys"] = float64(m.GCSys)
	c.gauges["HeapAlloc"] = float64(m.HeapAlloc)
	c.gauges["HeapIdle"] = float64(m.HeapIdle)
	c.gauges["HeapInuse"] = float64(m.HeapInuse)
	c.gauges["HeapObjects"] = float64(m.HeapObjects)
	c.gauges["HeapReleased"] = float64(m.HeapReleased)
	c.gauges["HeapSys"] = float64(m.HeapSys)
	c.gauges["LastGC"] = float64(m.LastGC)
	c.gauges["Lookups"] = float64(m.Lookups)
	c.gauges["MCacheInuse"] = float64(m.MCacheInuse)
	c.gauges["MCacheSys"] = float64(m.MCacheSys)
	c.gauges["MSpanInuse"] = float64(m.MSpanInuse)
	c.gauges["MSpanSys"] = float64(m.MSpanSys)
	c.gauges["Mallocs"] = float64(m.Mallocs)
	c.gauges["NextGC"] = float64(m.NextGC)
	c.gauges["NumForcedGC"] = float64(m.NumForcedGC)
	c.gauges["NumGC"] = float64(m.NumGC)
	c.gauges["OtherSys"] = float64(m.OtherSys)
	c.gauges["PauseTotalNs"] = float64(m.PauseTotalNs)
	c.gauges["StackInuse"] = float64(m.StackInuse)
	c.gauges["StackSys"] = float64(m.StackSys)
	c.gauges["Sys"] = float64(m.Sys)
	c.gauges["TotalAlloc"] = float64(m.TotalAlloc)

	// кастомная метрика RandomValue
	c.gauges["RandomValue"] = rand.Float64()

	// counter метрика PollCount (увеличивается при каждом вызове)
	c.counters["PollCount"]++

	fmt.Printf("[Collector] Собрано метрик: gauges=%d, counters=%d\n", len(c.gauges), len(c.counters))

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
