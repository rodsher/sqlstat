package sqlstat

import (
	"database/sql"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricOpenConnectionsTotal              = "open_connections_total"
	MetricConnectionsInUseTotal             = "connections_in_use_total"
	MetricConnectionsIdleTotal              = "connections_idle_total"
	MetricConnectionsWaitTotal              = "connections_wait_total"
	MetricConnectionsWaitDurationTotal      = "connections_wait_duration_total"
	MetricConnectionsMaxIdleClosedTotal     = "connections_max_idle_closed_total"
	MetricConnectionsMaxLifetimeClosedTotal = "connections_max_lifetime_closed_total"
	MetricMaxOpenConnections                = "max_open_connections"
)

// collectorCount describes how many Prometheus collectors exist
const collectorCount = 8

type (
	metric struct {
		Name  string
		Value int64
	}

	collector struct {
		DB *sql.DB
		collectorOpts
		metrics map[string]prometheus.Gauge
	}

	collectorOpts struct {
		Namespace string
		Subsystem string
		Interval  time.Duration
	}
)

func newCollector() *collector {
	return &collector{
		metrics: make(map[string]prometheus.Gauge),
	}
}

func (c *collector) withDB(db *sql.DB) *collector {
	c.DB = db

	return c
}

func (c *collector) withOpts(opts Opts) *collector {
	c.collectorOpts = collectorOpts(opts)

	return c
}

func (c *collector) registerMetrics() {
	c.registerOpenConnections()
	c.registerConnectionsInUse()
	c.registerConnectionsIdle()
	c.registerConnectionsWait()
	c.registerConnectionsWaitDuration()
	c.registerConnectionsMaxIdleClosed()
	c.registerConnectionsMaxLifetimeClosed()
	c.registerMaxOpenConnections()
}

func (c *collector) collectMetricsAsync() {
	go func() {
		var (
			metrics = make(chan metric)
			wg      = sync.WaitGroup{}
			delta   = 2
		)

		wg.Add(delta)

		go c.collectMetricsPeriodically(metrics)
		go c.updateMetrics(metrics)

		wg.Wait()
	}()
}

func (c *collector) collectMetricsPeriodically(metrics chan<- metric) {
	for {
		c.collectMetrics(metrics)
		time.Sleep(c.Interval)
	}
}

func (c *collector) collectMetrics(metrics chan<- metric) {
	metrics <- metric{MetricOpenConnectionsTotal, int64(c.DB.Stats().OpenConnections)}
	metrics <- metric{MetricConnectionsInUseTotal, int64(c.DB.Stats().InUse)}
	metrics <- metric{MetricConnectionsIdleTotal, int64(c.DB.Stats().Idle)}
	metrics <- metric{MetricConnectionsWaitTotal, c.DB.Stats().WaitCount}
	metrics <- metric{MetricConnectionsWaitDurationTotal, int64(c.DB.Stats().WaitDuration)}
	metrics <- metric{MetricConnectionsMaxIdleClosedTotal, c.DB.Stats().MaxIdleClosed}
	metrics <- metric{MetricConnectionsMaxLifetimeClosedTotal, c.DB.Stats().MaxLifetimeClosed}
	metrics <- metric{MetricMaxOpenConnections, int64(c.DB.Stats().MaxOpenConnections)}
}

func (c *collector) updateMetrics(metrics <-chan metric) {
	for m := range metrics {
		c.metrics[m.Name].Set(float64(m.Value))
	}
}

func (c *collector) registerOpenConnections() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      MetricOpenConnectionsTotal,
		Help:      "The number of established connections both in use and idle",
	})
	c.metrics[MetricOpenConnectionsTotal] = g
}

func (c *collector) registerConnectionsInUse() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      MetricConnectionsInUseTotal,
		Help:      "The number of connections currently in use",
	})
	c.metrics[MetricConnectionsInUseTotal] = g
}

func (c *collector) registerConnectionsIdle() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      MetricConnectionsIdleTotal,
		Help:      "The number of idle connections",
	})
	c.metrics[MetricConnectionsIdleTotal] = g
}

func (c *collector) registerConnectionsWait() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      MetricConnectionsWaitTotal,
		Help:      "The total number of connections waited for",
	})
	c.metrics[MetricConnectionsWaitTotal] = g
}

func (c *collector) registerConnectionsWaitDuration() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      MetricConnectionsWaitDurationTotal,
		Help:      "The total time blocked waiting for a new connection",
	})
	c.metrics[MetricConnectionsWaitDurationTotal] = g
}

func (c *collector) registerConnectionsMaxIdleClosed() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      MetricConnectionsMaxIdleClosedTotal,
		Help:      "The total number of connections closed due to SetMaxIdleConns",
	})
	c.metrics[MetricConnectionsMaxIdleClosedTotal] = g
}

func (c *collector) registerConnectionsMaxLifetimeClosed() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      MetricConnectionsMaxLifetimeClosedTotal,
		Help:      "The total number of connections closed due to SetConnMaxLifetime",
	})
	c.metrics[MetricConnectionsMaxLifetimeClosedTotal] = g
}

func (c *collector) registerMaxOpenConnections() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      MetricMaxOpenConnections,
		Help:      "Maximum number of open connections to the database",
	})
	c.metrics[MetricMaxOpenConnections] = g
}

func (c *collector) getCollectors() []prometheus.Collector {
	var collectors = make([]prometheus.Collector, 0, collectorCount)

	for _, metric := range c.metrics {
		collectors = append(collectors, metric)
	}

	return collectors
}
