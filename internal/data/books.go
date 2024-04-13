package data

import (
	"time"
)

type Book struct {
	ID        int64     `json:"id"` // Converts key name from ID to id
	CreatedAt time.Time `json:"-"`  // Removes this from response json
	Title     string    `json:"title"`
	Published int       `json:"published,omitempty"` // Makes Published optional
	Pages     int       `json:"pages,omitempty"`
	Geners    []string  `json:"geners,omitempty"`
	Raring    float32   `json:"rating,omitempty"`
	Version   int16     `json:"-"`
}
