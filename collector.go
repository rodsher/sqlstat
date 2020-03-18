package sqlstat

const (
	namespace = "sql"
	subsystem = "stats"
)

// Collector is a core concept of the package. Must be register in Prometheus after initilization.
type Collector struct {
	CollectorOpts
}

// CollectorOpts uses to overwrite default metrics properties.
//
// Example:
// 	opts := sqlstat.CollectorOpts{
//		Namespace: "project_name",
//		Subsystem: "database",
//	}
type CollectorOpts struct {
	Namespace string
	Subsystem string
}

// NewCollector creates ready to register collector with default or passed options.
//
// Default options:
//	c := sqlstat.NewCollector()
//	prometheus.MustRegister(c)
//
// Overwrite default options:
//	c := sqlstat.NewCollector(sqlstat.CollectorOpts{
//		Namespace: "project_name",
//		Subsystem: "database",
//	})
//	prometheus.MustRegister(c)
func NewCollector(options ...CollectorOpts) *Collector {
	opts := CollectorOpts{
		Namespace: namespace,
		Subsystem: subsystem,
	}

	if len(options) > 0 {
		opts = options[0]
	}

	return &Collector{
		CollectorOpts: opts,
	}
}
