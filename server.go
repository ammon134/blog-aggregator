package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
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

	config := &Config{
		Port: os.Getenv("PORT"),
	}

	svr := NewServer(config)

	httpServer := &http.Server{
		Addr:    ":" + config.Port,
		Handler: svr,
	}

	fmt.Printf("listening on port %s...\n", config.Port)
	err = httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
	return nil
}
