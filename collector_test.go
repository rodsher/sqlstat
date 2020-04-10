package sqlstat

import (
	"database/sql"
	"runtime"
	"testing"
	"time"
)

func TestNewCollector(t *testing.T) {
	c := newCollector()
	if c == nil {
		t.Error("must be initialized")
	}
}

func TestNewCollector_withoutPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error("unexpected panic", err)
		}
	}()

	newCollector()
}

func BenchmarkNewCollector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		newCollector()
	}
}

func BenchmarkNewCollector_WithOptsAndDB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		newCollector().withDB(&sql.DB{}).withOpts(Opts{
			Namespace: "sql",
			Subsystem: "stat",
			Interval:  5 * time.Second,
		})
	}
}

func TestNewCollector_WithDB(t *testing.T) {
	c := newCollector().withDB(&sql.DB{})
	if c.DB == nil {
		t.Error("must be initialized")
	}
}

func TestNewCollector_WithDB_withoutPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error("unexpected panic", err)
		}
	}()

	newCollector().withDB(&sql.DB{})
}

func BenchmarkNewCollector_WithDB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		newCollector().withDB(&sql.DB{})
	}
}

func TestNewCollector_WithOpts(t *testing.T) {
	c := newCollector().withOpts(Opts{
		Namespace: "sql",
		Subsystem: "stat",
		Interval:  5 * time.Second,
	})
	if c.Namespace != "sql" {
		t.Errorf("expect: %s, get: %s", "sql", c.Namespace)
	}

	if c.Subsystem != "stat" {
		t.Errorf("expect: %s, get: %s", "stat", c.Subsystem)
	}

	if c.Interval.Seconds() != 5.0 {
		t.Errorf("expect: %f, get: %f", 5.0, c.Interval.Seconds())
	}
}

func TestNewCollector_WithOpts_withoutPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error("unexpected panic")
		}
	}()

	newCollector().withOpts(Opts{
		Namespace: "sql",
		Subsystem: "stat",
		Interval:  5 * time.Second,
	})
}

func BenchmarkNewCollector_WithOpts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		newCollector().withOpts(Opts{
			Namespace: "sql",
			Subsystem: "stat",
			Interval:  5 * time.Second,
		})
	}
}

func TestCollector_RegisterMetrics(t *testing.T) {
	c := newCollector().withDB(&sql.DB{}).withOpts(Opts{
		Namespace: "sql",
		Subsystem: "stat",
		Interval:  5 * time.Second,
	})

	if len(c.metrics) != 0 {
		t.Errorf("expect: %d, get: %d", 0, len(c.metrics))
	}

	c.registerMetrics()

	if len(c.metrics) != collectorCount {
		t.Errorf("expect: %d, get: %d", collectorCount, len(c.metrics))
	}
}

func BenchmarkCollector_RegisterMetrics(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		db := sql.DB{}
		opts := Opts{
			Namespace: "sql",
			Subsystem: "stat",
			Interval:  5 * time.Second,
		}
		c := newCollector().withDB(&db).withOpts(opts)

		b.StartTimer()
		c.registerMetrics()
	}
}

func TestCollector_CollectMetricsAsync_numGoroutineIncreased(t *testing.T) {
	t.Parallel()

	var (
		n = runtime.NumGoroutine()
		c = newCollector().withDB(&sql.DB{}).withOpts(Opts{
			Namespace: "sql",
			Subsystem: "stat",
			Interval:  5 * time.Second,
		})
	)

	c.registerMetrics()
	c.collectMetricsAsync()

	// wait until runtime launched go routines
	time.Sleep(50 * time.Millisecond)

	if runtime.NumGoroutine() <= n {
		t.Error("num goroutine must be increased")
	}
}

func TestCollector_CollectMetricsPeriodically(t *testing.T) {
	metrics := make(chan metric)
	c := newCollector().withDB(&sql.DB{}).withOpts(Opts{
		Namespace: "sql",
		Subsystem: "stat",
		Interval:  5 * time.Second,
	})

	c.registerMetrics()
	go c.collectMetricsPeriodically(metrics)

	for i := 0; i < collectorCount; i++ {
		<-metrics
	}

	// If all messages were read from the metrics channel, metrics were collected
	// and written to the channel
	close(metrics)
}

func TestCollector_UpdateMetrics(t *testing.T) {
	metrics := make(chan metric)
	c := newCollector().withDB(&sql.DB{}).withOpts(Opts{
		Namespace: "sql",
		Subsystem: "stat",
		Interval:  5 * time.Second,
	})

	c.registerMetrics()
	go c.updateMetrics(metrics)

	for i := 0; i < collectorCount; i++ {
		metrics <- metric{"connections_in_use_total", 0}
	}

	// If loop completed, metrics were read from channel
	close(metrics)
}

func TestCollector_RegisterOpenConnections(t *testing.T) {
	c := newCollector().withDB(&sql.DB{}).withOpts(Opts{
		Namespace: "sql",
		Subsystem: "stat",
		Interval:  5 * time.Second,
	})

	c.registerOpenConnections()

	if len(c.metrics) != 1 {
		t.Errorf("expect: %d, get: %d", 1, len(c.metrics))
	}
}

func TestCollector_RegisterConnectionsInUse(t *testing.T) {
	c := newCollector().withDB(&sql.DB{}).withOpts(Opts{
		Namespace: "sql",
		Subsystem: "stat",
		Interval:  5 * time.Second,
	})

	c.registerConnectionsInUse()

	if len(c.metrics) != 1 {
		t.Errorf("expect: %d, get: %d", 1, len(c.metrics))
	}
}

func TestCollector_RegisterConnectionsIdle(t *testing.T) {
	c := newCollector().withDB(&sql.DB{}).withOpts(Opts{
		Namespace: "sql",
		Subsystem: "stat",
		Interval:  5 * time.Second,
	})

	c.registerConnectionsIdle()

	if len(c.metrics) != 1 {
		t.Errorf("expect: %d, get: %d", 1, len(c.metrics))
	}
}

func TestCollector_RegisterConnectionsWait(t *testing.T) {
	c := newCollector().withDB(&sql.DB{}).withOpts(Opts{
		Namespace: "sql",
		Subsystem: "stat",
		Interval:  5 * time.Second,
	})

	c.registerConnectionsWait()

	if len(c.metrics) != 1 {
		t.Errorf("expect: %d, get: %d", 1, len(c.metrics))
	}
}

func TestCollector_RegisterConnectionsWaitDuration(t *testing.T) {
	c := newCollector().withDB(&sql.DB{}).withOpts(Opts{
		Namespace: "sql",
		Subsystem: "stat",
		Interval:  5 * time.Second,
	})

	c.registerConnectionsWaitDuration()

	if len(c.metrics) != 1 {
		t.Errorf("expect: %d, get: %d", 1, len(c.metrics))
	}
}

func TestCollector_RegisterConnectionsMaxIdleClosed(t *testing.T) {
	c := newCollector().withDB(&sql.DB{}).withOpts(Opts{
		Namespace: "sql",
		Subsystem: "stat",
		Interval:  5 * time.Second,
	})

	c.registerConnectionsWaitDuration()

	if len(c.metrics) != 1 {
		t.Errorf("expect: %d, get: %d", 1, len(c.metrics))
	}
}

func TestCollector_RegisterConnectionsMaxLifetimeClosed(t *testing.T) {
	c := newCollector().withDB(&sql.DB{}).withOpts(Opts{
		Namespace: "sql",
		Subsystem: "stat",
		Interval:  5 * time.Second,
	})

	c.registerConnectionsMaxLifetimeClosed()

	if len(c.metrics) != 1 {
		t.Errorf("expect: %d, get: %d", 1, len(c.metrics))
	}
}

func TestCollector_RegisterMaxOpenConnections(t *testing.T) {
	c := newCollector().withDB(&sql.DB{}).withOpts(Opts{
		Namespace: "sql",
		Subsystem: "stat",
		Interval:  5 * time.Second,
	})

	c.registerMaxOpenConnections()

	if len(c.metrics) != 1 {
		t.Errorf("expect: %d, get: %d", 1, len(c.metrics))
	}
}
