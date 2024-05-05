package main

import (
	"context"
	"net/http"

	"github.com/ammon134/blog-aggregator/internal/auth"
)

type ContextKey string

const AuthDBUser ContextKey = "middleware.auth.DBUser"

func (c Config) middlewareIsAuthed(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apikey, err := auth.GetAuthToken(r.Header, auth.AuthTypeAPIKey)
		if err != nil {
			respondError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		user, err := c.DB.GetUserByApiKey(r.Context(), apikey)
		if err != nil {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}

		ctx := context.WithValue(r.Context(), AuthDBUser, &user)
		req := r.WithContext(ctx)
		h.ServeHTTP(w, req)
	})
}
