package main

import (
	"flag"
	"log"
	"net/http"
)

type application struct {
}

func main() {
	addr := flag.String("addr", ":80", "HTTP Network address")
	app := &application{}
	server := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	log.Printf("Starting the web application server %s", *addr)
	err := server.ListenAndServe()
	log.Fatal(err)
}
