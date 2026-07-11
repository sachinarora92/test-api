package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// zipCodeRegex is pre-compiled for efficient validation
var zipCodeRegex = regexp.MustCompile(`^[A-Za-z0-9\-]{3,10}$`)

// Address represents a physical address.
type Address struct {
	ID        string    `json:"id"`
	Street    string    `json:"street"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	Zip       string    `json:"zip"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ValidateAddress validates the Address fields.
func ValidateAddress(addr *Address) error {
	if addr.Street == "" {
		return fmt.Errorf("street is required")
	}
	if addr.City == "" {
		return fmt.Errorf("city is required")
	}
	if addr.Zip == "" {
		return fmt.Errorf("zip is required")
	}
	if addr.Country == "" {
		return fmt.Errorf("country is required")
	}

	// Validate zip format (basic: alphanumeric with optional hyphens, 3-10 chars)
	if !isValidZip(addr.Zip) {
		return fmt.Errorf("invalid zip format")
	}

	// State is optional (not required)

	return nil
}

// isValidZip validates basic zip/postal code format.
func isValidZip(zip string) bool {
	// Allow 3-10 alphanumeric characters with optional hyphens
	return zipCodeRegex.MatchString(zip)
}

// NewAddress creates a new Address with a generated ID and timestamps.
func NewAddress(street, city, state, zip, country string) *Address {
	now := time.Now().UTC()
	return &Address{
		ID:        uuid.New().String(),
		Street:    street,
		City:      city,
		State:     state,
		Zip:       zip,
		Country:   country,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
