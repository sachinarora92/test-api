package main

import (
	"testing"
)

func TestValidateAddress(t *testing.T) {
	tests := []struct {
		name      string
		address   *Address
		shouldErr bool
	}{
		{
			name: "valid address",
			address: &Address{
				Street:  "123 Main St",
				City:    "New York",
				State:   "NY",
				Zip:     "10001",
				Country: "USA",
			},
			shouldErr: false,
		},
		{
			name: "valid address without state",
			address: &Address{
				Street:  "123 Main St",
				City:    "London",
				Zip:     "SW1A-1AA",
				Country: "UK",
			},
			shouldErr: false,
		},
		{
			name: "missing street",
			address: &Address{
				City:    "New York",
				State:   "NY",
				Zip:     "10001",
				Country: "USA",
			},
			shouldErr: true,
		},
		{
			name: "missing city",
			address: &Address{
				Street:  "123 Main St",
				State:   "NY",
				Zip:     "10001",
				Country: "USA",
			},
			shouldErr: true,
		},
		{
			name: "missing zip",
			address: &Address{
				Street:  "123 Main St",
				City:    "New York",
				State:   "NY",
				Country: "USA",
			},
			shouldErr: true,
		},
		{
			name: "missing country",
			address: &Address{
				Street: "123 Main St",
				City:   "New York",
				State:  "NY",
				Zip:    "10001",
			},
			shouldErr: true,
		},
		{
			name: "invalid zip too short",
			address: &Address{
				Street:  "123 Main St",
				City:    "New York",
				State:   "NY",
				Zip:     "12",
				Country: "USA",
			},
			shouldErr: true,
		},
		{
			name: "invalid zip too long",
			address: &Address{
				Street:  "123 Main St",
				City:    "New York",
				State:   "NY",
				Zip:     "123456789012",
				Country: "USA",
			},
			shouldErr: true,
		},
		{
			name: "valid zip with hyphen",
			address: &Address{
				Street:  "123 Main St",
				City:    "London",
				Zip:     "SW1A-1AA",
				Country: "UK",
			},
			shouldErr: false,
		},
		{
			name: "invalid zip with special characters",
			address: &Address{
				Street:  "123 Main St",
				City:    "New York",
				State:   "NY",
				Zip:     "1000@#",
				Country: "USA",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAddress(tt.address)
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidateAddress() error = %v, shouldErr %v", err, tt.shouldErr)
			}
		})
	}
}

func TestIsValidZip(t *testing.T) {
	tests := []struct {
		name    string
		zip     string
		isValid bool
	}{
		{"US zip", "10001", true},
		{"US zip with hyphen", "10001-1234", true},
		{"UK postcode", "SW1A-1AA", true},
		{"Canadian postal code", "K1A-0B1", true},
		{"too short", "12", false},
		{"too long", "123456789012", false},
		{"with special characters", "1000@#", false},
		{"with spaces", "10001 1234", false},
		{"minimum length", "100", true},
		{"maximum length", "1234567890", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidZip(tt.zip)
			if got != tt.isValid {
				t.Errorf("isValidZip(%q) = %v, want %v", tt.zip, got, tt.isValid)
			}
		})
	}
}

func TestNewAddress(t *testing.T) {
	addr := NewAddress("123 Main St", "New York", "NY", "10001", "USA")

	if addr.ID == "" {
		t.Error("NewAddress should generate an ID")
	}

	if addr.Street != "123 Main St" {
		t.Errorf("expected street '123 Main St', got '%s'", addr.Street)
	}

	if addr.City != "New York" {
		t.Errorf("expected city 'New York', got '%s'", addr.City)
	}

	if addr.State != "NY" {
		t.Errorf("expected state 'NY', got '%s'", addr.State)
	}

	if addr.Zip != "10001" {
		t.Errorf("expected zip '10001', got '%s'", addr.Zip)
	}

	if addr.Country != "USA" {
		t.Errorf("expected country 'USA', got '%s'", addr.Country)
	}

	if addr.CreatedAt.IsZero() {
		t.Error("NewAddress should set CreatedAt")
	}

	if addr.UpdatedAt.IsZero() {
		t.Error("NewAddress should set UpdatedAt")
	}

	if !addr.CreatedAt.Equal(addr.UpdatedAt) {
		t.Error("NewAddress should set CreatedAt and UpdatedAt to the same time")
	}
}
