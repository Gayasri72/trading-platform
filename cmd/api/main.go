
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gayasri72/trading-platform/internal/database"
	"github.com/joho/godotenv"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func main() {
	// Load environment variables from .env (local development)
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found; using system environment variables")
	}

	// Database configuration
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// API port configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Connect to PostgreSQL
	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")

	// HTTP routes
	mux := http.NewServeMux()

	// Liveness check
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GoTrade API is alive"))
	})

	// Readiness check (includes database connectivity)
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		err := db.Ping(ctx)

		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(HealthResponse{
				Status:   "not_ready",
				Database: "disconnected",
			})
			return
		}

		json.NewEncoder(w).Encode(HealthResponse{
			Status:   "ready",
			Database: "connected",
		})
	})

	// HTTP server configuration
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("GoTrade API running on port %s", port)

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
