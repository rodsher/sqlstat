package sqlstat

import (
	"database/sql"
	"runtime"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestNew_defaultOpts(t *testing.T) {
	stat := New()
	if stat.GetOpts().Namespace != "sql" {
		t.Errorf("expect: %s, get: %s", "sql", stat.GetOpts().Namespace)
	}

	if stat.GetOpts().Subsystem != "stat" {
		t.Errorf("expect: %s, get: %s", "stat", stat.GetOpts().Subsystem)
	}

	if stat.GetOpts().Interval.Seconds() != 5.0 {
		t.Errorf("expect: %f, get: %f", 5.0, stat.GetOpts().Interval.Seconds())
	}
}

func TestNew_withOpts(t *testing.T) {
	stat := New(Opts{
		Namespace: "ns",
		Subsystem: "sb",
		Interval:  5 * time.Second,
	})
	if stat.GetOpts().Namespace != "ns" {
		t.Errorf("expect: %s, get: %s", "ns", stat.GetOpts().Namespace)
	}

	if stat.GetOpts().Subsystem != "sb" {
		t.Errorf("expect: %s, get: %s", "sb", stat.GetOpts().Subsystem)
	}

	if stat.GetOpts().Interval.Seconds() != 5.0 {
		t.Errorf("expect: %f, get: %f", 5.0, stat.GetOpts().Interval.Seconds())
	}
}

func TestNew_withoutPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error("unexpected panic")
		}
	}()

	New()
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if stat := New(); stat == nil {
			b.Error("must be initialized")
		}
	}
}

func BenchmarkOpts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opts := &Opts{
			Namespace: "ns",
			Subsystem: "sb",
			Interval:  5 * time.Second,
		}
		if opts == nil {
			b.Error("must be initialized")
		}
	}
}

func TestStat_RegisterDB(t *testing.T) {
	var (
		stat = New()
		db   = sql.DB{}
	)

	err := stat.RegisterDB(&db)
	if err != nil {
		t.Error("must not raise error when argument is initialized pointer")
	}
}

func TestStat_RegisterDB_allCollectorsEnabled(t *testing.T) {
	var (
		stat = New()
		db   = sql.DB{}
	)

	err := stat.RegisterDB(&db)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if len(stat.GetCollectors()) != 8 {
		t.Errorf("expect: %d, get: %d", 8, len(stat.GetCollectors()))
	}
}

func TestStat_RegisterDB_nilArgument(t *testing.T) {
	stat := New()

	err := stat.RegisterDB(nil)
	if err == nil {
		t.Error("must raise error when argument is nil")
	}
}

func TestStat_RegisterDB_withoutPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error("unexpected panic")
		}
	}()

	var (
		stat = New()
		db   = sql.DB{}
	)

	err := stat.RegisterDB(&db)
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestStat_RegisterDB_numGoroutineIncreased(t *testing.T) {
	var (
		stat = New()
		db   = sql.DB{}
		n    = runtime.NumGoroutine()
	)

	err := stat.RegisterDB(&db)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if runtime.NumGoroutine() < n {
		t.Error("num goroutine must be increased")
	}
}

func BenchmarkStat_RegisterDB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		var (
			stat = New()
			db   = sql.DB{}
		)

		b.StartTimer()
		//nolint:errcheck
		stat.RegisterDB(&db)
	}
}

func TestStat_GetCollectors(t *testing.T) {
	s := &stat{
		collectors: []prometheus.Collector{
			prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "max_open_connections",
			}),
		},
	}

	collectors := s.GetCollectors()
	if len(collectors) != 1 {
		t.Errorf("expect: %d, get: %d", 1, len(collectors))
	}
}
