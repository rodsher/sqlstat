package sqlstat

import "github.com/prometheus/client_golang/prometheus"

func (s *stat) enableCollectors() {
	s.enableOpenConnections()
	s.enableConnectionsInUse()
	s.enableConnectionsIdle()
	s.enableConnectionsWait()
	s.enableConnectionsWaitDuration()
	s.enableConnectionsMaxIdleClosed()
	s.enableConnectionsMaxLifetimeClosed()
	s.enableMaxOpenConnections()
}

func (s *stat) enableOpenConnections() {
	c := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: s.Namespace,
		Subsystem: s.Subsystem,
		Name:      "open_connections_total",
		Help:      "The number of established connections both in use and idle",
	})
	s.metrics["open_connections_total"] = c
	s.collectors = append(s.collectors, c)
}

func (s *stat) enableConnectionsInUse() {
	c := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: s.Namespace,
		Subsystem: s.Subsystem,
		Name:      "connections_in_use_total",
		Help:      "The number of connections currently in use",
	})
	s.collectors = append(s.collectors, c)
}

func (s *stat) enableConnectionsIdle() {
	c := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: s.Namespace,
		Subsystem: s.Subsystem,
		Name:      "connections_idle_total",
		Help:      "The number of idle connections",
	})
	s.collectors = append(s.collectors, c)
}

func (s *stat) enableConnectionsWait() {
	c := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: s.Namespace,
		Subsystem: s.Subsystem,
		Name:      "connections_wait_total",
		Help:      "The total number of connections waited for",
	})
	s.collectors = append(s.collectors, c)
}

func (s *stat) enableConnectionsWaitDuration() {
	c := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: s.Namespace,
		Subsystem: s.Subsystem,
		Name:      "connections_wait_duration_total",
		Help:      "The total time blocked waiting for a new connection",
	})
	s.collectors = append(s.collectors, c)
}

func (s *stat) enableConnectionsMaxIdleClosed() {
	c := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: s.Namespace,
		Subsystem: s.Subsystem,
		Name:      "connections_max_idle_closed_total",
		Help:      "The total number of connections closed due to SetMaxIdleConns",
	})
	s.collectors = append(s.collectors, c)
}

func (s *stat) enableConnectionsMaxLifetimeClosed() {
	c := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: s.Namespace,
		Subsystem: s.Subsystem,
		Name:      "connections_max_lifetime_closed_total",
		Help:      "The total number of connections closed due to SetConnMaxLifetime",
	})
	s.collectors = append(s.collectors, c)
}

func (s *stat) enableMaxOpenConnections() {
	c := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: s.Namespace,
		Subsystem: s.Subsystem,
		Name:      "max_open_connections",
		Help:      "Maximum number of open connections to the database",
	})
	s.collectors = append(s.collectors, c)
}
