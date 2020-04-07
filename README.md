# go-sql-prom
Client to collect connection metrics using database/sql

[![Build Status](https://travis-ci.org/rodsher/sqlstat.svg?branch=master)](https://travis-ci.org/rodsher/sqlstat)
[![Go Report Card](https://goreportcard.com/badge/github.com/rodsher/sqlstat)](https://goreportcard.com/report/github.com/rodsher/sqlstat)

## Exposed metrics
* sql_stats_open_connections_total - The number of established connections both in use and idle
* sql_stats_connections_in_use_total - The number of connections currently in use
* sql_stats_connections_idle_total - The number of idle connections
* sql_stats_connections_wait_total - The total number of connections waited for
* sql_stats_connections_wait_duration_total - The total time blocked waiting for a new connection
* sql_stats_connections_max_idle_closed_total - The total number of connections closed due to SetMaxIdleConns
* sql_stats_connections_max_lifetime_closed_total - The total number of connections closed due to SetConnMaxLifetime
* sql_stats_max_open_connections - Maximum number of open connections to the database


