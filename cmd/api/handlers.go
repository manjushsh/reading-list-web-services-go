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

	}
	if r.Method == http.MethodPost {
		var newBook struct {
			Title     string   `json:"title"`
			Published int      `json:"published"`
			Pages     int      `json:"pages"`
			Geners    []string `json:"geners"`
			Rating    float32  `json:"rating"`
		}

		err := app.readJSON(w, r, &newBook)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "%v\n", newBook)
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
		Rating:    5,
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
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	var newBook struct {
		Title     *string  `json:"title"`
		Published *int     `json:"published"`
		Pages     *int     `json:"pages"`
		Geners    []string `json:"geners"`
		Rating    *float32 `json:"rating"`
	}

	storedData := data.Book{
		ID:        idInt,
		CreatedAt: time.Now(),
		Title:     "Test",
		Published: 2019,
		Pages:     300,
		Geners:    []string{"f", "t"},
		Rating:    5,
		Version:   1,
	}

	err = app.readJSON(w, r, &newBook)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if newBook.Title != nil {
		storedData.Title = *newBook.Title
	}
	if newBook.Published != nil {
		storedData.Published = *newBook.Published
	}
	if newBook.Pages != nil {
		storedData.Pages = *newBook.Pages
	}
	if len(newBook.Geners) > 0 {
		storedData.Geners = newBook.Geners
	}
	if newBook.Rating != nil {
		storedData.Rating = *newBook.Rating
	}
	fmt.Fprintf(w, "%v\n", storedData)
}

func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
	}
	fmt.Fprintf(w, "Delete details of book with id %d", idInt)
}
