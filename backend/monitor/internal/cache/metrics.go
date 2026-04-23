package cache

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Metrics provides Prometheus counters and histograms for cache operations.
type Metrics struct {
	hits       prometheus.Counter
	misses     prometheus.Counter
	errors     *prometheus.CounterVec
	timeouts   prometheus.Counter
	latency    prometheus.Histogram
	writeBytes prometheus.Histogram
}

// NewMetrics registers and returns a new Metrics set.
func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		hits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_hit_total",
			Help: "Total number of cache hits",
		}),
		misses: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_miss_total",
			Help: "Total number of cache misses",
		}),
		errors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "cache_error_total",
			Help: "Total number of cache errors by type",
		}, []string{"error_type"}),
		timeouts: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_timeout_total",
			Help: "Total number of cache timeouts",
		}),
		latency: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "cache_latency_seconds",
			Help:    "Cache operation latency",
			Buckets: prometheus.DefBuckets,
		}),
		writeBytes: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "cache_write_bytes",
			Help:    "Size of cache writes in bytes",
			Buckets: prometheus.ExponentialBuckets(256, 4, 8),
		}),
	}

	reg.MustRegister(m.hits, m.misses, m.errors, m.timeouts, m.latency, m.writeBytes)
	return m
}

// RecordHit records a cache hit.
func (m *Metrics) RecordHit(_ string, d time.Duration) {
	m.hits.Inc()
	m.latency.Observe(d.Seconds())
}

// RecordMiss records a cache miss.
func (m *Metrics) RecordMiss(_ string, d time.Duration) {
	m.misses.Inc()
	m.latency.Observe(d.Seconds())
}

// RecordError records a cache error by type.
func (m *Metrics) RecordError(_ string, errType string) {
	m.errors.WithLabelValues(errType).Inc()
}

// RecordWrite records a cache write.
func (m *Metrics) RecordWrite(_ string, bytes int, _ time.Duration) {
	m.writeBytes.Observe(float64(bytes))
}
