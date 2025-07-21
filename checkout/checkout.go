package checkout

import (
	"fmt"

	"github.com/TheFodfather/checkoutapi/domain"
	"github.com/google/uuid"
)

// PricingService defines the dependency needed to get pricing rules.
type PricingService interface {
	GetRules() map[string]domain.PricingRule
}

type session struct {
	id           string
	scannedItems map[string]int
	pricer       PricingService // Dependency on the pricing service
}

// New creates a new checkout session instance.
func New(pricer PricingService) domain.ICheckout {
	return &session{
		id:           uuid.New().String(),
		scannedItems: make(map[string]int),
		pricer:       pricer,
	}
}

// Scan validates an SKU against the current pricing rules and adds it to the session.
func (s *session) Scan(SKU string) (err error) {
	rules := s.pricer.GetRules()
	if _, exists := rules[SKU]; !exists {
		return fmt.Errorf("sku '%s' not found in pricing rules", SKU)
	}
	s.scannedItems[SKU]++
	return nil
}

// GetTotalPrice calculates the total price for the session based on current pricing rules.
func (s *session) GetTotalPrice() (totalPrice int, err error) {
	rules := s.pricer.GetRules()
	totalPrice = 0
	for sku, count := range s.scannedItems {
		rule := rules[sku]
		if rule.SpecialPrice != nil && count >= rule.SpecialPrice.Quantity {
			numOffers := count / rule.SpecialPrice.Quantity
			totalPrice += numOffers * rule.SpecialPrice.Price
			remaining := count % rule.SpecialPrice.Quantity
			totalPrice += remaining * rule.UnitPrice
		} else {
			totalPrice += count * rule.UnitPrice
		}
	}
	return totalPrice, nil
}
