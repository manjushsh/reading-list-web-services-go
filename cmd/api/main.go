package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type application struct {
	config config
	logger *log.Logger
}

type config struct {
	port int
	env  string
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API Server Port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev | stage | prod)")
	flag.Parse()

	loggger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &application{
		config: cfg,
		logger: loggger,
	}

	addr := fmt.Sprintf(":%d", cfg.port)

	server := &http.Server{
		Addr:         addr,
		Handler:      app.route(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	loggger.Printf("Starting %s server on %s...", cfg.env, addr)
	err := server.ListenAndServe()
	loggger.Fatal(err)

}
