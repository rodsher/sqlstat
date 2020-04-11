# SQL stat
Package to collect client's connection metrics for database/sql.

[![Build Status](https://travis-ci.org/rodsher/sqlstat.svg?branch=master)](https://travis-ci.org/rodsher/sqlstat)
[![Coverage Status](https://coveralls.io/repos/github/rodsher/sqlstat/badge.svg?branch=master)](https://coveralls.io/github/rodsher/sqlstat?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rodsher/sqlstat)](https://goreportcard.com/report/github.com/rodsher/sqlstat)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/git-chglog/git-chglog/blob/master/LICENSE)

## Motivation

It is very important to monitor RDBMS metrics and performance.
But modern RDBMS can be used by multiple applications, and each application connected to the database using a pool of connections.
Collect connection pool statistics also very important to diagnose troubles in the system.
*sqlstat* provides a package to collect statistics on how your application uses the pool of connections.

## Prerequisites

Package requires a Go version with Modules support and uses versioned import. So please make sure to initialize a Go module before installing.

## Installation

```bash
go get github.com/rodsher/sqlstat
```

Import:

```go
import "github.com/rodsher/sqlstat"
```

## Quick start

Here we open connection with a PostgreSQL database server and register created connection pool
in sqlstat. Also, we register collectors in Prometheus.

```go
func main() {
	db, _ := sql.Open("postgres", "postgres://user:password@localhost/db")

	stat := sqlstat.New()
	stat.RegisterDB(db)

	prometheus.MustRegister(stat.GetCollectors()...)
}
```

Request Prometheus metrics at 8000 TCP port using cURL.

```bash
$ curl http://localhost:8000/metrics
# HELP sql_stat_connections_idle_total The number of idle connections
# TYPE sql_stat_connections_idle_total gauge
sql_stat_connections_idle_total 0
# HELP sql_stat_connections_in_use_total The number of connections currently in use
# TYPE sql_stat_connections_in_use_total gauge
sql_stat_connections_in_use_total 0
# HELP sql_stat_connections_max_idle_closed_total The total number of connections closed due to SetMaxIdleConns
# TYPE sql_stat_connections_max_idle_closed_total gauge
sql_stat_connections_max_idle_closed_total 0
# HELP sql_stat_connections_max_lifetime_closed_total The total number of connections closed due to SetConnMaxLifetime
# TYPE sql_stat_connections_max_lifetime_closed_total gauge
sql_stat_connections_max_lifetime_closed_total 0
# HELP sql_stat_connections_wait_duration_total The total time blocked waiting for a new connection
# TYPE sql_stat_connections_wait_duration_total gauge
sql_stat_connections_wait_duration_total 0
# HELP sql_stat_connections_wait_total The total number of connections waited for
# TYPE sql_stat_connections_wait_total gauge
sql_stat_connections_wait_total 0
# HELP sql_stat_max_open_connections Maximum number of open connections to the database
# TYPE sql_stat_max_open_connections gauge
sql_stat_max_open_connections 0
# HELP sql_stat_open_connections_total The number of established connections both in use and idle
# TYPE sql_stat_open_connections_total gauge
sql_stat_open_connections_total 0
```

## Exposed metrics

| Metric                                          | Description                                                 |
|-------------------------------------------------|-------------------------------------------------------------|
|   sql_stat_open_connections_total               |   The number of established connections both in use and idle  |
|   sql_stat_open_connections_total               |   The number of established connections both in use and idle|
|   sql_stat_connections_in_use_total             |   The number of connections currently in use|
|   sql_stat_connections_idle_total               |   The number of idle connections|
|   sql_stat_connections_wait_total               |   The total number of connections waited for|
|   sql_stat_connections_wait_duration_total      |   The total time blocked waiting for a new connection|
|   sql_stat_connections_max_idle_closed_total    |   The total number of connections closed due to SetMaxIdleConns|
|   sql_stat_connections_max_lifetime_closed_total|   The total number of connections closed due to SetConnMaxLifetime|
|   sql_stat_max_open_connections                 |   Maximum number of open connections to the database|

## PostgreSQL example

```go
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
```

Launch server

```bash
go run main.go
```

## Built with

[Prometheus](https://prometheus.io)

## Versioning

We use [SemVer](http://semver.org/) for versioning.

## License

This project is licensed under the MIT License - see the [LICENSE.md](https://github.com/rodsher/sqlstat/blob/master/LICENSE) file for details
