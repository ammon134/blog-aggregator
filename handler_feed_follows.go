package main

import (
	"net/http"
	"time"

	"github.com/ammon134/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func handleFeedFollowsCreate(config *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(AuthDBUser).(*database.User)

		type Parameters struct {
			FeedID uuid.UUID `json:"feed_id"`
		}

		params := &Parameters{}
		err := decode(r, params)
		if err != nil {
			respondError(w, http.StatusBadRequest, "missing parameters")
			return
		}

		dbFeedFollow, err := config.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    params.FeedID,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Class() == pq.ErrorClass("23") {
				respondError(w, http.StatusBadRequest, "feed_id is not valid")
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// NOTE: this only works because database.FeedFollow and FeedFollow
		// have the exact same structure
		feedFollow := FeedFollow(dbFeedFollow)

		respondJSON(w, http.StatusCreated, feedFollow)
	})
}
