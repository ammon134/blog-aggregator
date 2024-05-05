package main

import (
	"net/http"
	"time"

	"github.com/ammon134/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func handleFeedsCreate(config *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(AuthDBUser).(*database.User)

		type Parameters struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		}

		params := &Parameters{}
		err := decode(r, params)
		if err != nil {
			respondError(w, http.StatusBadRequest, "missing parameters")
			return
		}

		newFeedID := uuid.New()
		dbFeed, err := config.DB.CreateFeed(r.Context(), database.CreateFeedParams{
			ID:        newFeedID,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      params.Name,
			Url:       params.Url,
			UserID:    user.ID,
		})
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		dbFeedFollow, err := config.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    newFeedID,
		})
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// NOTE: this only works because database.Feed and Feed
		// have the exact same structure
		feed := Feed(dbFeed)
		feedFollow := FeedFollow(dbFeedFollow)

		type Response struct {
			Feed       Feed       `json:"feed"`
			FeedFollow FeedFollow `json:"feed_follow"`
		}

		respondJSON(w, http.StatusCreated, Response{
			Feed:       feed,
			FeedFollow: feedFollow,
		})
	})
}

func handleFeedsList(config *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbFeeds, err := config.DB.ListFeeds(r.Context())
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		feeds := []Feed{}
		for _, dbFeed := range dbFeeds {
			feeds = append(feeds, Feed(dbFeed))
		}

		respondJSON(w, http.StatusOK, feeds)
	})
}
