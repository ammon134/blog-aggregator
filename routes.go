package main

import (
	"net/http"
)

// addRoutes()
func addRoutes(mux *http.ServeMux, config *Config) {
	mux.Handle("GET /v1/readiness", handleReadiness())
	mux.Handle("GET /v1/err", handleError())

	mux.Handle("POST /v1/users", handleUsersCreate(config))
	mux.Handle("GET /v1/users", config.middlewareIsAuthed(handleUsersGet()))

}
