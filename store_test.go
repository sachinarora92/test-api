package main

import (
	"testing"
)

func TestStoreCreate(t *testing.T) {
	store := NewStore()

	addr := NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	created, err := store.Create(addr)

	if err != nil {
		t.Fatalf("failed to create address: %v", err)
	}

	if created.ID == "" {
		t.Error("expected address to have an ID")
	}

	if created.Street != addr.Street {
		t.Errorf("expected street %s, got %s", addr.Street, created.Street)
	}
}

func TestStoreCreateDuplicate(t *testing.T) {
	store := NewStore()

	addr := NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	store.Create(addr)

	// Try to create with same ID
	_, err := store.Create(addr)
	if err == nil {
		t.Error("expected error when creating duplicate address")
	}
}

func TestStoreGetByID(t *testing.T) {
	store := NewStore()

	addr := NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	created, _ := store.Create(addr)

	retrieved, err := store.GetByID(created.ID)
	if err != nil {
		t.Fatalf("failed to get address: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, retrieved.ID)
	}
}

func TestStoreGetByIDNotFound(t *testing.T) {
	store := NewStore()

	_, err := store.GetByID("non-existent-id")
	if err == nil {
		t.Error("expected error when getting non-existent address")
	}
}

func TestStoreList(t *testing.T) {
	store := NewStore()

	addr1 := NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	addr2 := NewAddress("456 Oak Ave", "Boston", "MA", "02101", "USA")

	store.Create(addr1)
	store.Create(addr2)

	addresses := store.List()
	if len(addresses) != 2 {
		t.Errorf("expected 2 addresses, got %d", len(addresses))
	}
}

func TestStoreListEmpty(t *testing.T) {
	store := NewStore()

	addresses := store.List()
	if len(addresses) != 0 {
		t.Errorf("expected 0 addresses, got %d", len(addresses))
	}
}

func TestStoreUpdate(t *testing.T) {
	store := NewStore()

	addr := NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	created, _ := store.Create(addr)

	updates := &Address{
		Street:  "456 New St",
		City:    "Boston",
		Zip:     "02101",
		Country: "USA",
	}

	updated, err := store.Update(created.ID, updates)
	if err != nil {
		t.Fatalf("failed to update address: %v", err)
	}

	if updated.Street != "456 New St" {
		t.Errorf("expected street to be updated to '456 New St', got '%s'", updated.Street)
	}

	if updated.City != "Boston" {
		t.Errorf("expected city to be updated to 'Boston', got '%s'", updated.City)
	}

	if updated.UpdatedAt.Before(updated.CreatedAt) {
		t.Error("expected UpdatedAt to be after or equal to CreatedAt")
	}
}

func TestStoreUpdateNotFound(t *testing.T) {
	store := NewStore()

	updates := &Address{Street: "New St"}
	_, err := store.Update("non-existent-id", updates)

	if err == nil {
		t.Error("expected error when updating non-existent address")
	}
}

func TestStoreDelete(t *testing.T) {
	store := NewStore()

	addr := NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	created, _ := store.Create(addr)

	err := store.Delete(created.ID)
	if err != nil {
		t.Fatalf("failed to delete address: %v", err)
	}

	_, err = store.GetByID(created.ID)
	if err == nil {
		t.Error("expected address to be deleted")
	}
}

func TestStoreDeleteNotFound(t *testing.T) {
	store := NewStore()

	err := store.Delete("non-existent-id")
	if err == nil {
		t.Error("expected error when deleting non-existent address")
	}
}
