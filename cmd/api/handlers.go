package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Status: %s\n", "Online")
	fmt.Fprintf(w, "Environment: %s \n ", app.config.env)
	fmt.Fprintf(w, "Version: %s\n", version)
}

func (app *application) getCreateBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Display the books on reading list.")
	}
	if r.Method == http.MethodPost {
		fmt.Fprintln(w, "Added a new book to reading list.")
	}
}

func (app *application) getUpdateDeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getBook(w, r)
	case http.MethodPut:
		app.updateBook(w, r)
	case http.MethodDelete:
		app.deleteBook(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
	}
	// file deepcode ignore XSS: <It is demo project so ignore it.>
	fmt.Fprintf(w, "Display details of book with id %d", idInt)
}

func (app *application) updateBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
	}
	fmt.Fprintf(w, "Update details of book with id %d", idInt)
}

func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
	}
	fmt.Fprintf(w, "Delete details of book with id %d", idInt)
}
