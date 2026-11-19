package order

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
	repo    *Repository
}

func NewHandler(service *Service, repo *Repository) *Handler {
	return &Handler{service: service, repo: repo}
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidOrder):
		respondJSON(w, 400, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrInsufficientFunds):
		respondJSON(w, 409, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrDuplicateOrder):
		respondJSON(w, 409, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		respondJSON(w, 404, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrNotCancellable):
		respondJSON(w, 409, map[string]string{"error": err.Error()})
	default:
		log.Printf("Order error: %v", err)
		respondJSON(w, 500, map[string]string{
			"error": "internal server error",
		})
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	var req CreateRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		respondJSON(w, 400, map[string]string{
			"error": "invalid JSON request",
		})
		return
	}

	// Reject trailing JSON values.
	var extra any
	if err := decoder.Decode(&extra); err == nil {
		respondJSON(w, 400, map[string]string{
			"error": "only one JSON object is allowed",
		})
		return
	} else if !errors.Is(err, io.EOF) {
		respondJSON(w, 400, map[string]string{
			"error": "invalid trailing JSON",
		})
		return
	}

	o, err := h.service.Create(r.Context(), req)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, o)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	orders, err := h.repo.List(r.Context())
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, 200, orders)
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		respondJSON(w, 400, map[string]string{
			"error": "invalid order ID",
		})
		return
	}

	if err := h.service.Cancel(r.Context(), id); err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, 200, map[string]string{
		"status": "CANCELLED",
	})
}
