package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"readinglist.manjushsh.github.io/internal/data"
)

const version = "1.0.0"

type application struct {
	config config
	logger *log.Logger
	models *data.Models
}

type config struct {
	port int
	env  string
	dsn  string // Data Name Service
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API Server Port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev | stage | prod)")
	flag.StringVar(&cfg.dsn, "db-dsn", os.Getenv("READINGLIST_DB_DSN"), "Postgres SQL DSN")
	flag.Parse()

	loggger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		loggger.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		loggger.Fatal(err)
	}
	loggger.Printf("Connected to db pool")

	addr := fmt.Sprintf(":%d", cfg.port)
	models := data.NewModels(db)
	app := &application{
		config: cfg,
		logger: loggger,
		models: &models,
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      app.route(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	loggger.Printf("Starting %s server on %s...", cfg.env, addr)
	err = server.ListenAndServe()
	loggger.Fatal(err)

}
