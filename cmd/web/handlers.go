package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Home Page")
}

func (app *application) bookView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "View Book")
}

func (app *application) bookCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create new Book form")
}
