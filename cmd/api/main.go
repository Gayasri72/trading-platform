package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gayasri72/trading-platform/internal/database"
	"github.com/gayasri72/trading-platform/internal/market"
	"github.com/gayasri72/trading-platform/internal/order"
	"github.com/gayasri72/trading-platform/internal/wallet"
	"github.com/joho/godotenv"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func main() {
	// Load environment variables
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

	// HTTP router
	mux := http.NewServeMux()

	// Health endpoints
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GoTrade API is alive"))
	})

	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(HealthResponse{
				Status:   "not_ready",
				Database: "disconnected",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{
			Status:   "ready",
			Database: "connected",
		})
	})

	// Day 2: Market module
	marketRepo := market.NewRepository(db)
	marketHandler := market.NewHandler(marketRepo)

	mux.HandleFunc("GET /api/v1/assets", marketHandler.ListAssets)
	mux.HandleFunc("GET /api/v1/assets/{symbol}", marketHandler.GetAsset)

	mux.HandleFunc("GET /api/v1/markets", marketHandler.ListMarkets)
	mux.HandleFunc("GET /api/v1/markets/{symbol}", marketHandler.GetMarket)
	//day 3
	walletRepo := wallet.NewRepository(db)
	walletHandler := wallet.NewHandler(walletRepo)

	mux.HandleFunc("GET /api/v1/wallets", walletHandler.GetBalances)
	mux.HandleFunc("GET /api/v1/wallets/ledger", walletHandler.GetLedger)

	//day 4
	orderRepo := order.NewRepository(db)
	orderService := order.NewService(orderRepo)
	orderHandler := order.NewHandler(orderService, orderRepo)

	mux.HandleFunc("POST /api/v1/orders", orderHandler.Create)
	mux.HandleFunc("GET /api/v1/orders", orderHandler.List)
	mux.HandleFunc("DELETE /api/v1/orders/{id}", orderHandler.Cancel)

	// HTTP server
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
