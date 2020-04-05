package sqlstat

import (
	"database/sql"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

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

func (c *collector) withDB(DB *sql.DB) *collector {
	c.DB = DB

	return c
}

func (c *collector) withOpts(opts Opts) *collector {
	c.collectorOpts = collectorOpts{
		Namespace: opts.Namespace,
		Subsystem: opts.Subsystem,
		Interval:  opts.Interval,
	}

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
		)

		wg.Add(2)

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
	metrics <- metric{"open_connections_total", int64(c.DB.Stats().OpenConnections)}
	metrics <- metric{"connections_in_use_total", int64(c.DB.Stats().InUse)}
	metrics <- metric{"connections_idle_total", int64(c.DB.Stats().Idle)}
	metrics <- metric{"connections_wait_total", int64(c.DB.Stats().WaitCount)}
	metrics <- metric{"connections_wait_duration_total", int64(c.DB.Stats().WaitDuration)}
	metrics <- metric{"connections_max_idle_closed_total", int64(c.DB.Stats().MaxIdleClosed)}
	metrics <- metric{"connections_max_lifetime_closed_total", int64(c.DB.Stats().MaxLifetimeClosed)}
	metrics <- metric{"max_open_connections", int64(c.DB.Stats().MaxOpenConnections)}
}

func (c *collector) updateMetrics(metrics <-chan metric) {
	for {
		select {
		case m := <-metrics:
			c.metrics[m.Name].Set(float64(m.Value))
		}
	}
}

func (c *collector) registerOpenConnections() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      "open_connections_total",
		Help:      "The number of established connections both in use and idle",
	})
	c.metrics["open_connections_total"] = g
}

func (c *collector) registerConnectionsInUse() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      "connections_in_use_total",
		Help:      "The number of connections currently in use",
	})
	c.metrics["connections_in_use_total"] = g
}

func (c *collector) registerConnectionsIdle() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      "connections_idle_total",
		Help:      "The number of idle connections",
	})
	c.metrics["connections_idle_total"] = g
}

func (c *collector) registerConnectionsWait() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      "connections_wait_total",
		Help:      "The total number of connections waited for",
	})
	c.metrics["connections_wait_total"] = g
}

func (c *collector) registerConnectionsWaitDuration() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      "connections_wait_duration_total",
		Help:      "The total time blocked waiting for a new connection",
	})
	c.metrics["connections_wait_duration_total"] = g
}

func (c *collector) registerConnectionsMaxIdleClosed() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      "connections_max_idle_closed_total",
		Help:      "The total number of connections closed due to SetMaxIdleConns",
	})
	c.metrics["connections_max_idle_closed_total"] = g
}

func (c *collector) registerConnectionsMaxLifetimeClosed() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      "connections_max_lifetime_closed_total",
		Help:      "The total number of connections closed due to SetConnMaxLifetime",
	})
	c.metrics["connections_max_lifetime_closed_total"] = g
}

func (c *collector) registerMaxOpenConnections() {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      "max_open_connections",
		Help:      "Maximum number of open connections to the database",
	})
	c.metrics["max_open_connections"] = g
}

func (c *collector) getCollectors() []prometheus.Collector {
	var collectors []prometheus.Collector

	for _, metric := range c.metrics {
		collectors = append(collectors, metric)
	}

	return collectors
}
