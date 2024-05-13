package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ammon134/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

const (
	QueryDefaultLimit = 20
)

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func handlePostsListByUserId(c *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(AuthDBUser).(*database.User)

		params := &database.ListPostsByUserIDParams{
			UserID: user.ID,
			Limit:  int32(QueryDefaultLimit),
		}
		limitStr := r.URL.Query().Get("limit")
		if limitStr != "" {
			l64, err := strconv.Atoi(limitStr)
			if err == nil {
				params.Limit = int32(l64)
			}
		}

		posts, err := c.DB.ListPostsByUserID(r.Context(), *params)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, createResponsePostList(posts))
	})
}

// Helper functions
func createResponsePost(dbPost database.Post) Post {
	description := ""
	if dbPost.Description.Valid {
		description = dbPost.Description.String
	}
	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Url:         dbPost.Url,
		Description: description,
		PublishedAt: dbPost.PublishedAt,
		FeedID:      dbPost.FeedID,
	}
}

func createResponsePostList(dbPosts []database.Post) []Post {
	posts := make([]Post, len(dbPosts))
	for i, p := range dbPosts {
		posts[i] = createResponsePost(p)
	}
	return posts
}
