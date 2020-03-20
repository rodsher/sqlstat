package sqlstat

import (
	"database/sql"
	"testing"
)

func TestNew_defaultOpts(t *testing.T) {
	stat := New()
	if stat.Namespace != "sql" {
		t.Errorf("expect: %s, get: %s", "sql", stat.Namespace)
	}

	if stat.Subsystem != "stat" {
		t.Errorf("expect: %s, get: %s", "stat", stat.Subsystem)
	}

	if stat.IsEnabled != true {
		t.Errorf("expect: %t, get: %t", true, stat.IsAllEnabled)
	}
}

func TestNew_withOpts(t *testing.T) {
	stat := New(StatOpts{
		Namespace: "ns",
		Subsystem: "sb",
		IsEnabled: false,
	})
	if stat.Namespace != "ns" {
		t.Errorf("expect: %s, get: %s", "ns", stat.Namespace)
	}

	if stat.Subsystem != "sb" {
		t.Errorf("expect: %s, get: %s", "sb", stat.Subsystem)
	}

	if stat.IsEnabled != false {
		t.Errorf("expect: %t, get: %t", false, stat.IsEnabled)
	}
}

func TestNew_withoutPanic() {
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

func BenchmarkStatOpts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opts := &StatOpts{
			Namespace: "ns",
			Subsystem: "sb",
			IsEnabled: false,
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
