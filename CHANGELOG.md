
<a name="v1.0.0.."></a>
## [v1.0.0..](https://github.com/rodsher/sqlstat/compare/v1.0.0...v1.0.0..) (2020-04-03)

### Docs

* **src:** Add docs for a register interface

### Feat

* **ci:** Create confuguration file for CI
* **examples:** Add example with PostgreSQL
* **repository:** Create .gitattributes file
* **repository:** Init go modules
* **repository:** Add Dockerfile for a golangci-lint and changelog config
* **scripts:** Create script for a generating CHANGELOG file
* **src:** Create scaffold for a stat concept
* **src:** Do new API more consise
* **src:** Add docs, constants, prettify API
* **src:** Add enablers method for a stat structure
* **src:** Add collector with max idle closed
* **src:** Create Collector contstructor and options
* **src:** Add methods to collect metrics
* **src:** Remove unused dependency lib pq
* **src:** Add error ErrDatabaseIsNil

### Fix

* **src:** Fix compile time error for collector.go file
* **src:** Fix broken tests for Collector

### Refactor

* **src:** Move collector to separate source file
* **src:** Remove collector concept from library

### Test

* **src:** Add unit tests and benchmarks for a stat
* **src:** Add unit tests and benchmarks for a collector
* **src:** Add unit test for collectors check
* **src:** Add unit tests for enablers method
* **src:** Improve code coverage, add tests for GetCollectors()
* **src:** Repair broken unit tests
* **src:** Add benchmark and acceptance test for Enable method
* **src:** Add acceptance tests for a stat.go
* **src:** Add high-level tests for Collector
* **src:** Create acceptance tests for a collector


<a name="v1.0.0"></a>
## v1.0.0 (2020-03-18)

### Docs

* **readme:** Add list of collected metrics

