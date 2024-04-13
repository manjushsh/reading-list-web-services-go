package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"readinglist.manjushsh.github.io/internal/data"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	data := map[string]string{
		"status":      "Online",
		"environment": app.config.env,
		"version":     version,
	}
	if err := app.writeJSON(w, http.StatusOK, envolope{"data": data}); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
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

	book := data.Book{
		ID:        idInt,
		CreatedAt: time.Now(),
		Title:     "Test",
		Published: 2019,
		Pages:     300,
		Geners:    []string{"f", "t"},
		Raring:    5,
		Version:   1,
	}

	if err := app.writeJSON(w, http.StatusOK, envolope{"book": book}); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
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
