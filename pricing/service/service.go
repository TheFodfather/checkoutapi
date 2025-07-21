package pricing

import (
	"sync"

	"github.com/TheFodfather/checkoutapi/domain"
)

// Service provides access to pricing rules
type Service struct {
	pricingFile string
	rules       map[string]domain.PricingRule
	sync.RWMutex
}
