package data

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Book struct {
	ID        int64     `json:"id"` // Converts key name from ID to id
	CreatedAt time.Time `json:"-"`  // Removes this from response json
	Title     string    `json:"title"`
	Published int       `json:"published,omitempty"` // Makes Published optional
	Pages     int       `json:"pages,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Rating    float32   `json:"rating,omitempty"`
	Version   int16     `json:"-"`
}

type BookModel struct {
	DB *sql.DB
}

func (b BookModel) InsertBook(book *Book) error {
	query := `
	INSERT INTO books (title, published, pages, genres, rating)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, version`
	args := []interface{}{book.Title, book.Published, book.Pages, pq.Array(book.Genres), book.Rating}
	// return the auto generated system values to Go object
	return b.DB.QueryRow(query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}

func (b BookModel) GetBook(id int64) (*Book, error) {
	if id < 1 {
		return nil, errors.New("record not found")
	}
	query := `
	SELECT id, title, published, pages, genres, rating, created_at, version
	FROM books
	WHERE id = $1`

	var book Book
	var genresDB []byte // Assuming the genres column type is []byte
	err := b.DB.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Published,
		&book.Pages,
		&genresDB,
		&book.Rating,
		&book.CreatedAt,
		&book.Version,
	)
	// Unmarshal genresDB into a []string
	book.Genres = strings.Split(string(genresDB), ",")

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("no entries found for given id")
		default:
			return nil, err
		}
	}
	return &book, nil
}

func (b BookModel) UpdateBook(book *Book) error {
	query := `
	UPDATE books
	SET title = $1, published = $2, pages = $3, genres = $4, rating = $5, version = version + 1
	WHERE id = $6 AND version = $7
	RETURNING version`
	args := []interface{}{book.Title, book.Published, book.Pages, pq.Array(book.Genres), book.Rating, book.ID, book.Version}
	return b.DB.QueryRow(query, args...).Scan(&book.Version)
}

func (b BookModel) DeleteBook(id int64) error {
	if id < 1 {
		return errors.New("record not found")
	}
	query := `
	DELETE FROM books
	WHERE id = $1`

	results, err := b.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}

func (b BookModel) GetAllBooks() ([]*Book, error) {
	query := `
	SELECT id, title, published, pages, genres, rating, created_at, version
	FROM books
	ORDER BY id`

	rows, err := b.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []*Book{}
	for rows.Next() {
		var book Book
		var genresDB []byte
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Published,
			&book.Pages,
			&genresDB,
			&book.Rating,
			&book.CreatedAt,
			&book.Version,
		)
		if err != nil {
			return nil, err
		}
		book.Genres = strings.Split(string(genresDB), ",")
		books = append(books, &book)
	}

	if rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}
