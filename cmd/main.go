package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rodsher/sqlstat"
)

var pool *sql.DB

func main() {
	pool, err := sql.Open("postgres", "postgres://safe_portal:safe_portal@localhost/safe_portal?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	pool.SetMaxOpenConns(16)

	defer func() {
		if err := pool.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := pool.QueryContext(ctx, "SELECT * FROM \"user\";")
	if err != nil {
		log.Fatal(err)
	}

	users := []struct {
		id    int64
		email string
	}{}
	for rows.Next() {
		user := struct {
			id    int64
			email string
		}{}
		err = rows.Scan(&user.id, &user.email)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	for _, user := range users {
		fmt.Println(user.id, user.email)
	}

	stat := sqlstat.New()
	err = stat.RegisterDB(pool)
	if err != nil {
		log.Fatal(err)
	}

	prometheus.MustRegister(stat.GetCollectors()...)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
