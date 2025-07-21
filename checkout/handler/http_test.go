package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TheFodfather/checkoutapi/checkout/repository"
	"github.com/TheFodfather/checkoutapi/domain"
)

// mockHandlerPricingService provides a mock implementation of the PricingService for testing.
type mockHandlerPricingService struct{}

// GetRules returns a predefined set of pricing rules for testing purposes.
func (m *mockHandlerPricingService) GetRules() map[string]domain.PricingRule {
	return map[string]domain.PricingRule{
		"A": {UnitPrice: 50, SpecialPrice: &domain.SpecialPrice{Quantity: 3, Price: 130}},
		"B": {UnitPrice: 30, SpecialPrice: &domain.SpecialPrice{Quantity: 2, Price: 45}},
		"C": {UnitPrice: 20},
		"D": {UnitPrice: 15},
	}
}

func TestCreateCheckout(t *testing.T) {
	server := setupTestServer(t)

	t.Run("create a new checkout session successfully", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/checkouts", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var body map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
			t.Fatalf("Could not parse response body: %v", err)
		}

		if _, ok := body["checkoutId"]; !ok {
			t.Error("Response body should contain a checkoutId")
		}
	})

	t.Run("create unique checkout IDs for multiple requests", func(t *testing.T) {
		id1 := createCheckoutSession(t, server)
		id2 := createCheckoutSession(t, server)

		if id1 == "" || id2 == "" {
			t.Fatal("Failed to create one or more checkout sessions")
		}

		if id1 == id2 {
			t.Errorf("Expected unique checkout IDs, but got the same ID twice: %s", id1)
		}
	})
}

func TestScanItem(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should scan a valid item successfully", func(t *testing.T) {
		checkoutID := createCheckoutSession(t, server)
		scanURL := "/checkouts/" + checkoutID + "/scan"
		payload := []byte(`{"sku":"A"}`)

		req, _ := http.NewRequest("POST", scanURL, bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}
	})

	t.Run("return 400 Bad Request for an invalid SKU", func(t *testing.T) {
		checkoutID := createCheckoutSession(t, server)
		scanURL := "/checkouts/" + checkoutID + "/scan"
		payload := []byte(`{"sku":"Z"}`) // Z is not a valid SKU

		req, _ := http.NewRequest("POST", scanURL, bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("return 404 Not Found for a non-existent checkout session", func(t *testing.T) {
		scanURL := "/checkouts/non-existent-id/scan"
		payload := []byte(`{"sku":"A"}`)

		req, _ := http.NewRequest("POST", scanURL, bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})
}

func TestGetTotalPrice(t *testing.T) {
	server := setupTestServer(t)

	t.Run("zero total price for a new checkout", func(t *testing.T) {
		checkoutID := createCheckoutSession(t, server)
		getURL := "/checkouts/" + checkoutID

		req, _ := http.NewRequest("GET", getURL, nil)
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var body struct {
			TotalPrice int `json:"totalPrice"`
		}
		if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
			t.Fatalf("Could not parse response body: %v", err)
		}

		if body.TotalPrice != 0 {
			t.Errorf("Expected total price to be 0 for a new checkout, got %d", body.TotalPrice)
		}
	})

	t.Run("correct total price after scanning items", func(t *testing.T) {
		checkoutID := createCheckoutSession(t, server)
		scanURL := "/checkouts/" + checkoutID + "/scan"

		// Scan A, B, A -> Total = 50 + 30 + 50 = 130
		itemsToScan := []string{"A", "B", "A"}
		expectedTotal := 130

		for _, sku := range itemsToScan {
			payload := []byte(`{"sku":"` + sku + `"}`)
			req, _ := http.NewRequest("POST", scanURL, bytes.NewBuffer(payload))
			rr := httptest.NewRecorder()
			server.ServeHTTP(rr, req)
			if status := rr.Code; status != http.StatusNoContent {
				t.Fatalf("Failed to scan item %s: got status %v", sku, status)
			}
		}

		getURL := "/checkouts/" + checkoutID
		req, _ := http.NewRequest("GET", getURL, nil)
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var body struct {
			TotalPrice int `json:"totalPrice"`
		}
		if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
			t.Fatalf("Could not parse response body: %v", err)
		}

		if body.TotalPrice != expectedTotal {
			t.Errorf("Incorrect total price: got %d want %d", body.TotalPrice, expectedTotal)
		}
	})

	t.Run("return 404 Not Found for a non-existent checkout session", func(t *testing.T) {
		getURL := "/checkouts/non-existent-id"
		req, _ := http.NewRequest("GET", getURL, nil)
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})
}

// setupTestServer initializes a new test server with an in-memory repository and a mock pricer.
func setupTestServer(t *testing.T) http.Handler {
	repo := repository.NewInMemoryRepository()
	pricer := &mockHandlerPricingService{}
	handler := New(repo, pricer)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux
}

func createCheckoutSession(t *testing.T, server http.Handler) string {
	req, _ := http.NewRequest("POST", "/checkouts", nil)
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Fatalf("Failed to create checkout session: got status %v, want %v", status, http.StatusCreated)
	}

	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("Could not parse response body for new checkout: %v", err)
	}

	checkoutID, ok := body["checkoutId"]
	if !ok || checkoutID == "" {
		t.Fatal("Response body does not contain a valid checkoutId")
	}
	return checkoutID
}
