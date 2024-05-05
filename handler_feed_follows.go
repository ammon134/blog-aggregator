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

		type parameters struct {
			FeedID uuid.UUID `json:"feed_id"`
		}

		params := &parameters{}
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
				respondError(w, http.StatusBadRequest, "invalid feed_id format")
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// NOTE: this only works because database.FeedFollow and FeedFollow
		// have the exact same structure.
		feedFollow := FeedFollow(dbFeedFollow)

		respondJSON(w, http.StatusCreated, feedFollow)
	})
}

func handleFeedFollowsDelete(config *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feedFollowIDStr := r.PathValue("feedFollowID")
		feedFollowID := &uuid.UUID{}
		err := feedFollowID.UnmarshalText([]byte(feedFollowIDStr))
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid feed_follow_id format")
			return
		}

		user := r.Context().Value(AuthDBUser).(*database.User)

		err = config.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
			ID:     *feedFollowID,
			UserID: user.ID,
		})
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, http.StatusText(http.StatusOK))
	})
}

func handleFeedFollowsListByUserID(config *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(AuthDBUser).(*database.User)
		dbFeedFollows, err := config.DB.ListFeedFollowsByUserId(r.Context(), user.ID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		feedFollows := make([]FeedFollow, len(dbFeedFollows))
		for i, f := range dbFeedFollows {
			feedFollows[i] = FeedFollow(f)
		}

		respondJSON(w, http.StatusOK, feedFollows)
	})
}
