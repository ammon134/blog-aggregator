package main

import "net/http"

// addRoutes()
func addRoutes(mux *http.ServeMux, config *Config) {
	mux.Handle("/", http.NotFoundHandler())
	_ = config
}
