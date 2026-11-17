package market

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) ListAssets(w http.ResponseWriter, r *http.Request) {
	assets, err := h.repo.ListAssets(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve assets", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, assets)
}

func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) {
	symbol := strings.ToUpper(r.PathValue("symbol"))

	asset, err := h.repo.GetAsset(r.Context(), symbol)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to retrieve asset", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, asset)
}

func (h *Handler) ListMarkets(w http.ResponseWriter, r *http.Request) {
	markets, err := h.repo.ListMarkets(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve markets", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, markets)
}

func (h *Handler) GetMarket(w http.ResponseWriter, r *http.Request) {
	symbol := strings.ToUpper(r.PathValue("symbol"))

	item, err := h.repo.GetMarket(r.Context(), symbol)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Market not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to retrieve market", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, item)
}
