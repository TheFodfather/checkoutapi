package checkout

import (
	"fmt"
	"testing"

	"github.com/TheFodfather/checkoutapi/domain"
)

type mockPricingService struct{}

func (m *mockPricingService) GetRules() map[string]domain.PricingRule {
	return map[string]domain.PricingRule{
		"A": {UnitPrice: 50, SpecialPrice: &domain.SpecialOffer{Quantity: 3, Price: 130}},
		"B": {UnitPrice: 30, SpecialPrice: &domain.SpecialOffer{Quantity: 2, Price: 45}},
		"C": {UnitPrice: 20, SpecialPrice: nil},
		"D": {UnitPrice: 15, SpecialPrice: nil},
	}
}

func TestScan(t *testing.T) {
	mockPricer := &mockPricingService{}

	testCases := []struct {
		name          string
		initialItems  map[string]int
		skuToScan     string
		expectedErr   error
		expectedItems map[string]int
	}{
		{
			name:         "Scan existing SKU 'A' successfully",
			initialItems: make(map[string]int),
			skuToScan:    "A",
			expectedErr:  nil,
			expectedItems: map[string]int{
				"A": 1,
			},
		},
		{
			name:         "Scan existing SKU 'B' multiple times",
			initialItems: map[string]int{"B": 2},
			skuToScan:    "B",
			expectedErr:  nil,
			expectedItems: map[string]int{
				"B": 3,
			},
		},
		{
			name:         "Scan existing SKU 'C' once",
			initialItems: make(map[string]int),
			skuToScan:    "C",
			expectedErr:  nil,
			expectedItems: map[string]int{
				"C": 1,
			},
		},
		{
			name:          "Scan non-existent SKU 'Z'",
			initialItems:  make(map[string]int),
			skuToScan:     "Z",
			expectedErr:   fmt.Errorf("sku 'Z' not found in pricing rules"),
			expectedItems: make(map[string]int),
		},
		{
			name:         "Scan with initial items and then a new SKU 'D'",
			initialItems: map[string]int{"A": 1, "B": 1},
			skuToScan:    "D",
			expectedErr:  nil,
			expectedItems: map[string]int{
				"A": 1,
				"B": 1,
				"D": 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := &session{
				id:           "test-session-id",
				scannedItems: tc.initialItems,
				pricer:       mockPricer,
			}

			err := s.Scan(tc.skuToScan)

			//check for expected error
			if tc.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected error: %v, but got no error", tc.expectedErr)
				} else if err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error message: %q, but got: %q", tc.expectedErr.Error(), err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}

			// Check if scannedItems map is updated as expected
			if len(s.scannedItems) != len(tc.expectedItems) {
				t.Errorf("Expected scannedItems length %d, got %d", len(tc.expectedItems), len(s.scannedItems))
			}
			for sku, count := range tc.expectedItems {
				if s.scannedItems[sku] != count {
					t.Errorf("For SKU %q: Expected count %d, got %d", sku, count, s.scannedItems[sku])
				}
			}
		})
	}
}

func TestCheckoutTotalPrice(t *testing.T) {
	testCases := []struct {
		name          string
		skusToScan    []string
		expectedTotal int
		expectedErr   bool
	}{
		{name: "Empty Cart", skusToScan: []string{}, expectedTotal: 0},
		{name: "Simple Items, No Offers", skusToScan: []string{"A", "B", "C"}, expectedTotal: 100},
		{name: "Special Offer Exact Match", skusToScan: []string{"A", "A", "A"}, expectedTotal: 130},
		{name: "Special Offer With Remainder", skusToScan: []string{"A", "A", "A", "A"}, expectedTotal: 180},
		{name: "Multiple Different Special Offers", skusToScan: []string{"A", "B", "A", "B", "A"}, expectedTotal: 175},
		{name: "Items Scanned Out of Order", skusToScan: []string{"B", "A", "B"}, expectedTotal: 95},
		{name: "Invalid SKU Scan", skusToScan: []string{"A", "Z"}, expectedErr: true},
		{name: "Comprehensive Mix", skusToScan: []string{"C", "B", "A", "B", "A", "A", "D"}, expectedTotal: 210},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPricer := &mockPricingService{}
			co := New(mockPricer)
			var scanErr error
			for _, sku := range tc.skusToScan {
				if err := co.Scan(sku); err != nil {
					scanErr = err
					break
				}
			}
			if tc.expectedErr {
				if scanErr == nil {
					t.Errorf("Expected an error but got none")
				}
				return
			}
			if scanErr != nil {
				t.Fatalf("Got unexpected error during scan: %v", scanErr)
			}
			total, err := co.GetTotalPrice()
			if err != nil {
				t.Fatalf("GetTotalPrice() returned an unexpected error: %v", err)
			}
			if total != tc.expectedTotal {
				t.Errorf("Expected total price to be %d, but got %d", tc.expectedTotal, total)
			}
		})
	}
}
