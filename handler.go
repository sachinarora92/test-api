package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler handles HTTP requests for the address API.
type Handler struct {
	store *Store
}

// NewHandler creates a new handler with the given store.
func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

// ErrorResponse is a standard error response format.
type ErrorResponse struct {
	Error string `json:"error"`
}

// ListResponse is the response format for list operations.
type ListResponse struct {
	Addresses []*Address `json:"addresses"`
	Count     int        `json:"count"`
}

// CreateAddress handles POST /addresses.
func (h *Handler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	var input Address

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	// Create new address with generated ID and timestamps
	addr := NewAddress(input.Street, input.City, input.State, input.Zip, input.Country)

	created, err := h.store.Create(addr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetAddress handles GET /addresses/{id}.
func (h *Handler) GetAddress(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	addr, err := h.store.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(addr)
}

// ListAddresses handles GET /addresses.
func (h *Handler) ListAddresses(w http.ResponseWriter, r *http.Request) {
	addresses := h.store.List()

	response := ListResponse{
		Addresses: addresses,
		Count:     len(addresses),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateAddress handles PUT /addresses/{id}.
func (h *Handler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input Address

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	updated, err := h.store.Update(id, &input)
	if err != nil {
		// Distinguish between "not found" and validation errors
		if err.Error() == fmt.Sprintf("address with id %s not found", id) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			// Validation error
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

// DeleteAddress handles DELETE /addresses/{id}.
func (h *Handler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.store.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// writeError writes a standard error response.
func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// Health is a simple health check endpoint.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
