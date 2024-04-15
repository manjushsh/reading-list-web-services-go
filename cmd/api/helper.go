package main

import (
	"encoding/json"
	"errors"
	"io"
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	// Limit max body size
	maxBodySize := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBodySize))

	// Decode
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return err
	}

	err := decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body should have only one json object")
	}
	return nil
}
