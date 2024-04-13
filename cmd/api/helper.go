package main

import (
	"encoding/json"
	"net/http"
)

type envolope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data any) error {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	jsonData = append(jsonData, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
	return nil
}
