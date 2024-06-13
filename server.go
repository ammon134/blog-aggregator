package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ammon134/blog-aggregator/internal/database"
	"github.com/joho/godotenv"
)

type Config struct {
	DB   *database.Queries
	Port string
}

// NewServer()
func NewServer(config *Config) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		config,
	)

	var handler http.Handler = mux
	return handler
}

// run()
func run() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("port is missing in env")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is missing in env")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error opening db connection")
	}

	dbQueries := database.New(db)

	config := &Config{
		Port: port,
		DB:   dbQueries,
	}

	svr := NewServer(config)

	httpServer := &http.Server{
		Addr:    ":" + config.Port,
		Handler: svr,
	}

	go runScrapingWorker(config, 10, time.Second*60)

	fmt.Printf("listening on port %s...\n", config.Port)
	err = httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
	return nil
}
