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

	mux.Handle("POST /v1/feeds", config.middlewareIsAuthed(handleFeedsCreate(config)))
	mux.Handle("GET /v1/feeds", handleFeedsList(config))

	mux.Handle("POST /v1/feed_follows", config.middlewareIsAuthed(handleFeedFollowsCreate(config)))
}
