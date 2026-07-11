package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestCreateAddress(t *testing.T) {
	store := NewStore()
	handler := NewHandler(store)
	router := setupRouter(handler)

	tests := []struct {
		name           string
		input          map[string]string
		expectedStatus int
		shouldHaveID   bool
	}{
		{
			name: "valid address",
			input: map[string]string{
				"street":  "123 Main St",
				"city":    "New York",
				"state":   "NY",
				"zip":     "10001",
				"country": "USA",
			},
			expectedStatus: http.StatusCreated,
			shouldHaveID:   true,
		},
		{
			name: "valid address without state",
			input: map[string]string{
				"street":  "123 Main St",
				"city":    "London",
				"zip":     "SW1A-1AA",
				"country": "UK",
			},
			expectedStatus: http.StatusCreated,
			shouldHaveID:   true,
		},
		{
			name: "missing street",
			input: map[string]string{
				"city":    "New York",
				"state":   "NY",
				"zip":     "10001",
				"country": "USA",
			},
			expectedStatus: http.StatusBadRequest,
			shouldHaveID:   false,
		},
		{
			name: "missing city",
			input: map[string]string{
				"street":  "123 Main St",
				"state":   "NY",
				"zip":     "10001",
				"country": "USA",
			},
			expectedStatus: http.StatusBadRequest,
			shouldHaveID:   false,
		},
		{
			name: "missing zip",
			input: map[string]string{
				"street":  "123 Main St",
				"city":    "New York",
				"state":   "NY",
				"country": "USA",
			},
			expectedStatus: http.StatusBadRequest,
			shouldHaveID:   false,
		},
		{
			name: "invalid zip format",
			input: map[string]string{
				"street":  "123 Main St",
				"city":    "New York",
				"state":   "NY",
				"zip":     "12",
				"country": "USA",
			},
			expectedStatus: http.StatusBadRequest,
			shouldHaveID:   false,
		},
		{
			name: "missing country",
			input: map[string]string{
				"street": "123 Main St",
				"city":   "New York",
				"state":  "NY",
				"zip":    "10001",
			},
			expectedStatus: http.StatusBadRequest,
			shouldHaveID:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/addresses", bytes.NewReader(body))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.shouldHaveID {
				var addr Address
				json.NewDecoder(w.Body).Decode(&addr)
				if addr.ID == "" {
					t.Error("expected address to have an ID")
				}
				if addr.CreatedAt.IsZero() {
					t.Error("expected address to have created_at timestamp")
				}
			}
		})
	}
}

func TestGetAddress(t *testing.T) {
	store := NewStore()
	handler := NewHandler(store)
	router := setupRouter(handler)

	// Create a test address
	addr := NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	store.Create(addr)

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		shouldExist    bool
	}{
		{
			name:           "existing address",
			id:             addr.ID,
			expectedStatus: http.StatusOK,
			shouldExist:    true,
		},
		{
			name:           "non-existent address",
			id:             "non-existent-id",
			expectedStatus: http.StatusNotFound,
			shouldExist:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/addresses/"+tt.id, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.shouldExist {
				var retrieved Address
				json.NewDecoder(w.Body).Decode(&retrieved)
				if retrieved.ID != addr.ID {
					t.Errorf("expected ID %s, got %s", addr.ID, retrieved.ID)
				}
			}
		})
	}
}

func TestListAddresses(t *testing.T) {
	store := NewStore()
	handler := NewHandler(store)
	router := setupRouter(handler)

	// Create test addresses
	addr1 := NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	addr2 := NewAddress("456 Oak Ave", "Boston", "MA", "02101", "USA")
	store.Create(addr1)
	store.Create(addr2)

	req := httptest.NewRequest("GET", "/addresses", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response ListResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Count != 2 {
		t.Errorf("expected 2 addresses, got %d", response.Count)
	}

	if len(response.Addresses) != 2 {
		t.Errorf("expected 2 addresses in array, got %d", len(response.Addresses))
	}
}

func TestUpdateAddress(t *testing.T) {
	store := NewStore()
	handler := NewHandler(store)
	router := setupRouter(handler)

	// Create a test address
	addr := NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	store.Create(addr)

	tests := []struct {
		name           string
		id             string
		update         map[string]string
		expectedStatus int
		shouldUpdate   bool
	}{
		{
			name: "valid update",
			id:   addr.ID,
			update: map[string]string{
				"street": "456 Oak Ave",
				"city":   "Boston",
				"zip":    "02101",
			},
			expectedStatus: http.StatusOK,
			shouldUpdate:   true,
		},
		{
			name: "invalid zip format",
			id:   addr.ID,
			update: map[string]string{
				"zip": "12",
			},
			expectedStatus: http.StatusBadRequest,
			shouldUpdate:   false,
		},
		{
			name:           "non-existent address",
			id:             "non-existent-id",
			update:         map[string]string{"street": "New St"},
			expectedStatus: http.StatusNotFound,
			shouldUpdate:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.update)
			req := httptest.NewRequest("PUT", "/addresses/"+tt.id, bytes.NewReader(body))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.shouldUpdate {
				var updated Address
				json.NewDecoder(w.Body).Decode(&updated)
				if updated.Street != tt.update["street"] {
					t.Errorf("expected street %s, got %s", tt.update["street"], updated.Street)
				}
			}
		})
	}
}

func TestDeleteAddress(t *testing.T) {
	store := NewStore()
	handler := NewHandler(store)
	router := setupRouter(handler)

	// Create a test address
	addr := NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	store.Create(addr)

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		shouldDelete   bool
	}{
		{
			name:           "existing address",
			id:             addr.ID,
			expectedStatus: http.StatusNoContent,
			shouldDelete:   true,
		},
		{
			name:           "non-existent address",
			id:             "non-existent-id",
			expectedStatus: http.StatusNotFound,
			shouldDelete:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/addresses/"+tt.id, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.shouldDelete {
				// Verify it's deleted by trying to get it
				req := httptest.NewRequest("GET", "/addresses/"+tt.id, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				if w.Code != http.StatusNotFound {
					t.Error("expected address to be deleted")
				}
			}
		})
	}
}

func TestHealth(t *testing.T) {
	store := NewStore()
	handler := NewHandler(store)
	router := setupRouter(handler)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)

	if response["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", response["status"])
	}
}

func TestConcurrentOperations(t *testing.T) {
	store := NewStore()

	// Create multiple addresses concurrently
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func(index int) {
			addr := NewAddress("Street "+string(rune(index)), "City", "ST", "10001", "USA")
			store.Create(addr)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}

	// Verify all addresses were created
	addresses := store.List()
	if len(addresses) != 5 {
		t.Errorf("expected 5 addresses, got %d", len(addresses))
	}
}
