package main

import (
	"net/http"
	"time"

	"github.com/ammon134/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type Feed struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Name        string     `json:"name"`
	Url         string     `json:"url"`
	UserID      uuid.UUID  `json:"user_id"`
	LastFetchAt *time.Time `json:"last_fetch_at"`
}

func handleFeedsCreate(c *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(AuthDBUser).(*database.User)

		type parameters struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		}

		params := &parameters{}
		err := decode(r, params)
		if err != nil {
			respondError(w, http.StatusBadRequest, "missing parameters")
			return
		}

		dbFeed, err := c.DB.CreateFeed(r.Context(), database.CreateFeedParams{
			ID:        uuid.New(),
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

		dbFeedFollow, err := c.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    dbFeed.ID,
		})
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// NOTE: this only works because database.Feed and Feed
		// have the exact same structure. Create helper func otherwise.
		// See handler_users.go
		feed := createResponseFeed(dbFeed)
		feedFollow := FeedFollow(dbFeedFollow)

		type response struct {
			Feed       Feed       `json:"feed"`
			FeedFollow FeedFollow `json:"feed_follow"`
		}

		respondJSON(w, http.StatusCreated, response{
			Feed:       feed,
			FeedFollow: feedFollow,
		})
	})
}

func handleFeedsList(c *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbFeeds, err := c.DB.ListFeeds(r.Context())
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		feeds := createResponseFeedList(dbFeeds)
		respondJSON(w, http.StatusOK, feeds)
	})
}

// Helper functions
func createResponseFeed(df database.Feed) Feed {
	return Feed{
		ID:        df.ID,
		CreatedAt: df.CreatedAt,
		UpdatedAt: df.UpdatedAt,
		Name:      df.Name,
		Url:       df.Url,
		UserID:    df.UserID,
	}
}

func createResponseFeedList(dfs []database.Feed) []Feed {
	feeds := make([]Feed, len(dfs))
	for i, dbFeed := range dfs {
		feeds[i] = createResponseFeed(dbFeed)
	}
	return feeds
}
