package wallet

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetBalances(w http.ResponseWriter, r *http.Request) {
	balances, err := h.repo.GetBalances(r.Context(), 1)
	if err != nil {
		log.Printf("GetBalances: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balances)
}

func (h *Handler) GetLedger(w http.ResponseWriter, r *http.Request) {
	entries, err := h.repo.GetLedger(r.Context(), 1)
	if err != nil {
		log.Printf("GetLedger: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}
