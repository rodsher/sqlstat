package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rodsher/sqlstat"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:password@localhost/db")
	if err != nil {
		log.Fatal(err)
	}

	stat := sqlstat.New()
	err = stat.RegisterDB(db)
	if err != nil {
		log.Fatal(err)
	}

	prometheus.MustRegister(stat.GetCollectors()...)

	http.ListenAndServe(":8000", promhttp.Handler())
}
