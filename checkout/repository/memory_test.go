package repository

import (
	"testing"

	"github.com/google/uuid"
)

type mockCheckout struct{ id string }

func (m *mockCheckout) Scan(SKU string) (err error)                { return nil }
func (m *mockCheckout) GetTotalPrice() (totalPrice int, err error) { return 0, nil }
func (m *mockCheckout) GetID() string                              { return m.id }
func (m *mockCheckout) GetScannedItems() map[string]int            { return nil }

func TestGetNotFound(t *testing.T) {
	repo := NewInMemoryRepository()

	// Test case: Attempt to get a nonexistent session
	_, err := repo.Get("nonexistent-id")
	if err == nil {
		t.Fatal("Expected error for getting nonexistent session, but got nil")
	}
}

func TestSaveAndGet(t *testing.T) {
	repo := NewInMemoryRepository()
	testID := uuid.New().String()
	co := &mockCheckout{id: testID}

	// Test case: Save a session and then retrieve it
	if err := repo.Save(co); err != nil {
		t.Fatalf("Save() returned an unexpected error: %v", err)
	}

	retrieved, err := repo.Get(testID)
	if err != nil {
		t.Fatalf("Get() returned an unexpected error: %v", err)
	}

	if retrieved.GetID() != testID {
		t.Errorf("Expected retrieved ID to be %s, got %s", testID, retrieved.GetID())
	}
}
