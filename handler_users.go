package main

import (
	"net/http"
	"time"

	"github.com/ammon134/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type User struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
	ID        uuid.UUID `json:"id"`
}

func handleUsersCreate(c *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Name string
		}
		params := &parameters{}
		err := decode(r, params)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := c.DB.CreateUser(r.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      params.Name,
		})
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		type response struct {
			User User `json:"user"`
		}
		respondJSON(w, http.StatusOK, response{
			User: createResponseUser(user),
		})
	})
}

func handleUsersGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(AuthDBUser).(*database.User)
		if !ok {
			respondError(w, http.StatusBadRequest, "user not found")
			return
		}

		type response struct {
			User `json:"user"`
		}
		respondJSON(w, http.StatusOK, response{
			User: createResponseUser(*user),
		})
	})
}

// Helper functions
func createResponseUser(u database.User) User {
	return User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.CreatedAt,
		Name:      u.Name,
		APIKey:    u.ApiKey,
	}
}
