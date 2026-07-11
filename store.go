package main

import (
	"fmt"
	"sync"
	"time"
)

// Store provides in-memory persistence for addresses.
type Store struct {
	mu        sync.RWMutex
	addresses map[string]*Address
}

// NewStore creates a new in-memory address store.
func NewStore() *Store {
	return &Store{
		addresses: make(map[string]*Address),
	}
}

// Create adds a new address to the store.
func (s *Store) Create(addr *Address) (*Address, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := ValidateAddress(addr); err != nil {
		return nil, err
	}

	if _, exists := s.addresses[addr.ID]; exists {
		return nil, fmt.Errorf("address with id %s already exists", addr.ID)
	}

	s.addresses[addr.ID] = addr
	return addr, nil
}

// GetByID retrieves an address by its ID.
func (s *Store) GetByID(id string) (*Address, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	addr, exists := s.addresses[id]
	if !exists {
		return nil, fmt.Errorf("address with id %s not found", id)
	}

	return addr, nil
}

// List retrieves all addresses.
func (s *Store) List() []*Address {
	s.mu.RLock()
	defer s.mu.RUnlock()

	addresses := make([]*Address, 0, len(s.addresses))
	for _, addr := range s.addresses {
		addresses = append(addresses, addr)
	}

	return addresses
}

// Update updates an existing address.
func (s *Store) Update(id string, updates *Address) (*Address, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	addr, exists := s.addresses[id]
	if !exists {
		return nil, fmt.Errorf("address with id %s not found", id)
	}

	// Create a working copy to validate before modifying original
	working := *addr

	// Apply updates to working copy
	if updates.Street != "" {
		working.Street = updates.Street
	}
	if updates.City != "" {
		working.City = updates.City
	}
	if updates.State != "" {
		working.State = updates.State
	}
	if updates.Zip != "" {
		working.Zip = updates.Zip
	}
	if updates.Country != "" {
		working.Country = updates.Country
	}

	// Validate the updated address (original unchanged if validation fails)
	if err := ValidateAddress(&working); err != nil {
		return nil, err
	}

	working.UpdatedAt = time.Now().UTC()
	s.addresses[id] = &working
	return &working, nil
}

// Delete removes an address from the store.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.addresses[id]
	if !exists {
		return fmt.Errorf("address with id %s not found", id)
	}

	delete(s.addresses, id)
	return nil
}
