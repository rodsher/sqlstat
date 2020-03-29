package sqlstat

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestNewCollector(t *testing.T) {
	c := newCollector()
	if c == nil {
		t.Error("must be initialized")
	}
}

func TestCollector_RegisterMetrics(t *testing.T) {
	c := &collector{}

	if len(c.metrics) != 0 {
		t.Errorf("expect: %d, get: %d", 0, len(c.metrics))
	}

	c.registerMetrics()

	if len(c.metrics) != 8 {
		t.Errorf("expect: %d, get: %d", 8, len(c.metrics))
	}
}

func TestCollect(t *testing.T) {
	metrics := make(chan<- map[string]prometheus.Metric)
}
