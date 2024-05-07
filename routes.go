package main

import (
	"net/http"
)

// addRoutes()
func addRoutes(mux *http.ServeMux, c *Config) {
	mux.Handle("GET /v1/readiness", handleReadiness())
	mux.Handle("GET /v1/err", handleError())

	mux.Handle("POST /v1/users", handleUsersCreate(c))
	mux.Handle("GET /v1/users", c.middlewareIsAuthed(handleUsersGet()))

	mux.Handle("POST /v1/feeds", c.middlewareIsAuthed(handleFeedsCreate(c)))
	mux.Handle("GET /v1/feeds", handleFeedsList(c))

	mux.Handle("POST /v1/feed_follows", c.middlewareIsAuthed(handleFeedFollowsCreate(c)))
	mux.Handle("DELETE /v1/feed_follows/{feedFollowID}", c.middlewareIsAuthed(handleFeedFollowsDelete(c)))
	mux.Handle("GET /v1/feed_follows", c.middlewareIsAuthed(handleFeedFollowsListByUserID(c)))

	mux.Handle("GET /v1/posts", c.middlewareIsAuthed(handlePostsListByUserId(c)))
}
