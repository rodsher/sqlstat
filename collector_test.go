package sqlstat

import (
	"testing"
)

func TestNewCollector_withoutPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error("unexpected panic", err)
		}
	}()

	c := NewCollector()
	if c == nil {
		t.Error("must be initialized structure")
	}
}

func TestNewCollector_defaultOpts(t *testing.T) {
	c := NewCollector()
	if c.Namespace != "sql" {
		t.Errorf("expect: %s, get: %s", "sql", c.Namespace)
	}

	if c.Subsystem != "stats" {
		t.Errorf("expect: %s, get: %s", "stats", c.Subsystem)
	}
}

func TestNewCollector_withOpts(t *testing.T) {
	c := NewCollector(CollectorOpts{
		Namespace: "ns",
		Subsystem: "subsystem",
	})
	if c.Namespace != "ns" {
		t.Errorf("expect: %s, get: %s", "ns", c.Namespace)
	}

	if c.Subsystem != "subsystem" {
		t.Errorf("expect: %s, get: %s", "subsystem", c.Subsystem)
	}
}

func BenchmarkNewCollector_defaultOpts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if с := NewController(); c == nil {
			b.Fatal("must be pointer")
		}
	}
}

func BenchmarkNewCollector_withOpts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if c := NewCollector(); c == nil {
			b.Fatal("must be pointer")
		}
	}
}
