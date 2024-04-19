package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	if err := app.writeJSON(w, http.StatusOK, envolope{"data": data}, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) getCreateBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		books, err := app.models.Books.GetAllBooks()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := app.writeJSON(w, http.StatusOK, envolope{"books": books}, nil); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if r.Method == http.MethodPost {
		var newBook struct {
			Title     string   `json:"title"`
			Published int      `json:"published"`
			Pages     int      `json:"pages"`
			Genres    []string `json:"genres"`
			Rating    float32  `json:"rating"`
		}

		err := app.readJSON(w, r, &newBook)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		book := &data.Book{
			Title:     newBook.Title,
			Published: newBook.Published,
			Pages:     newBook.Pages,
			Genres:    newBook.Genres,
			Rating:    newBook.Rating,
		}

		err = app.models.Books.InsertBook(book)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("v1/books/%d", book.ID))

		err = app.writeJSON(w, http.StatusCreated, envolope{"book": book}, headers)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// fmt.Fprintf(w, "%v\n", book)
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

	// book := data.Book{
	// 	ID:        idInt,
	// 	CreatedAt: time.Now(),
	// 	Title:     "Test",
	// 	Published: 2019,
	// 	Pages:     300,
	// 	Genres:    []string{"f", "t"},
	// 	Rating:    5,
	// 	Version:   1,
	// }

	book, err := app.models.Books.GetBook(idInt)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("record not found")):
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envolope{"book": book}, nil); err != nil {
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

	bookFromDB, err := app.models.Books.GetBook(idInt)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("record not found")):
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	var newBook struct {
		Title     *string  `json:"title"`
		Published *int     `json:"published"`
		Pages     *int     `json:"pages"`
		Genres    []string `json:"genres"`
		Rating    *float32 `json:"rating"`
	}

	// storedData := data.Book{
	// 	ID:        idInt,
	// 	CreatedAt: time.Now(),
	// 	Title:     "Test",
	// 	Published: 2019,
	// 	Pages:     300,
	// 	Genres:    []string{"f", "t"},
	// 	Rating:    5,
	// 	Version:   1,
	// }

	err = app.readJSON(w, r, &newBook)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if newBook.Title != nil {
		bookFromDB.Title = *newBook.Title
	}
	if newBook.Published != nil {
		bookFromDB.Published = *newBook.Published
	}
	if newBook.Pages != nil {
		bookFromDB.Pages = *newBook.Pages
	}
	if len(newBook.Genres) > 0 {
		bookFromDB.Genres = newBook.Genres
	}
	if newBook.Rating != nil {
		bookFromDB.Rating = *newBook.Rating
	}

	err = app.models.Books.UpdateBook(bookFromDB)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envolope{"book": bookFromDB}, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// fmt.Fprintf(w, "%v\n", bookFromDB)
}

func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
	}

	err = app.models.Books.DeleteBook(idInt)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("record not found")):
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envolope{"message": "Book removed"}, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Delete details of book with id %d", idInt)
}
