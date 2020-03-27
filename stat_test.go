package sqlstat

import (
	"database/sql"
	"testing"

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

	if stat.GetOpts().IsStatEnable != true {
		t.Errorf("expect: %t, get: %t", true, stat.GetOpts().IsStatEnable)
	}
}

func TestNew_withOpts(t *testing.T) {
	stat := New(Opts{
		Namespace:    "ns",
		Subsystem:    "sb",
		IsStatEnable: false,
	})
	if stat.GetOpts().Namespace != "ns" {
		t.Errorf("expect: %s, get: %s", "ns", stat.GetOpts().Namespace)
	}

	if stat.GetOpts().Subsystem != "sb" {
		t.Errorf("expect: %s, get: %s", "sb", stat.GetOpts().Subsystem)
	}

	if stat.GetOpts().IsStatEnable != false {
		t.Errorf("expect: %t, get: %t", false, stat.GetOpts().IsStatEnable)
	}
}

func TestNew_withoutPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error("unexpected panic")
		}
	}()

	stat := New()
	if stat == nil {
		t.Error("must be initialized")
	}
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
			Namespace:    "ns",
			Subsystem:    "sb",
			IsStatEnable: false,
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

func BenchmarkStat_RegisterDB(b *testing.B) {
	var (
		stat = New()
		db   = sql.DB{}
	)

	for i := 0; i < b.N; i++ {
		if err := stat.RegisterDB(&db); err != nil {
			b.Error("unexpected error", err)
		}
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

func TestStat_GetCollectors_empty(t *testing.T) {
	s := New()
	collectors := s.GetCollectors()
	if len(collectors) != 0 {
		t.Errorf("expect: %d, get: %d", 0, len(collectors))
	}
}
