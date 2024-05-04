package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondJSON[T any](w http.ResponseWriter, code int, v T) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("error encoding JSON: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondError(w http.ResponseWriter, code int, msg string) {
	type Err struct {
		Error string `json:"error"`
	}
	respondJSON(w, code, Err{
		Error: msg,
	})
}
