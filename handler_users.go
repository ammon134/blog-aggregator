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
	ID        uuid.UUID `json:"id"`
}

func handleUsersCreate(config *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Parameters struct {
			Name string
		}
		params := &Parameters{}
		err := decode(r, params)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := config.DB.CreateUser(r.Context(), database.CreateUserParams{
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
			User: User{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Name:      user.Name,
			},
		})
	})
}
