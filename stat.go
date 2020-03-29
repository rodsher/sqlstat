package sqlstat

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Namespace stores a default value for Prometheus collector namespace.
	Namespace = "sql"

	// Subsystem stores a default value for Prometheus collector subsystem.
	Subsystem = "stat"
)

// Stat is a main concept. Stat provides top-level functionality to register a database
// connection and register metrics in Prometheus.
//
// Example:
// 		db, err := sql.Open("postgres", "postgres://user:password@localhost/db")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		stat := sqlstat.New()
// 		err = stat.RegisterDB(&db)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		prometheus.MustRegister(stat.GetCollectors()...)
type Stat interface {
	register
	collectorsKeeper
	optsKeeper
}

// New creates ready to use statistic with default options or with an overwritten
// options.
//
// Example:
// 		stat := New()
func New(options ...Opts) Stat {
	opts := Opts{
		Namespace:    Namespace,
		Subsystem:    Subsystem,
		IsStatEnable: true,
	}

	if len(options) > 0 {
		opts = options[0]
	}

	return &stat{
		Opts:       opts,
		metrics:    make(map[string]prometheus.Metric),
		collectors: make([]prometheus.Collector, 8),
	}
}

// Opts describes settings for collectors.
//
// Example:
// 		opts := &Opts{
//			Namespace: "custom_namespace",
//			Subsystem: "custom_subsystem",
//			IsStatEnable: true,
// 		}
type Opts struct {
	Namespace    string
	Subsystem    string
	IsStatEnable bool
}

type register interface {
	// RegisterDB registers a database connection and enables
	// all collectors for metrics.
	// You must register only one database per sqlstat instance.
	// Method returns error when passed argument is nil.
	//
	// Example:
	//		stat := sqlstat.New()
	// 		err = stat.RegisterDB(&db)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	RegisterDB(*sql.DB) error
}

type collectorsKeeper interface {
	// GetCollectors aims to access registered Prometheus collectors
	// and register them in Prometheus.
	//
	// Example:
	// 		stat := sqlstat.New()
	// 		err = stat.RegisterDB(&db)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	//
	//		collectors := stat.GetCollectors()
	// 		prometheus.MustRegister(collectors...)
	GetCollectors() []prometheus.Collector
}

type optsKeeper interface {
	// GetOpts aims to access stat options for logging, debugging, etc.
	//
	// Example:
	// 		opts := stat.New().GetOpts()
	//		fmt.Println(opts.Namespace) // Output: "sql"
	//		fmt.Println(opts.Subsystem) // Output: "stat"
	//		fmt.Println(opts.IsStatEnabled) // Output: true
	GetOpts() *Opts
}

type stat struct {
	Opts
	DB         *sql.DB
	metrics    map[string]prometheus.Metric
	collectors []prometheus.Collector
}

func (s *stat) RegisterDB(db *sql.DB) error {
	return nil
}

func (s *stat) GetCollectors() []prometheus.Collector {
	return s.collectors
}

func (s *stat) GetOpts() *Opts {
	return &s.Opts
}
