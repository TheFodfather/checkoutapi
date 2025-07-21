package domain

// ICheckout is the core interface defining the contract for a checkout
type ICheckout interface {
	Scan(SKU string) (err error)
	GetTotalPrice() (totalPrice int, err error)
}

// PricingRule defines the pricing structure for a single SKU.
type PricingRule struct {
	UnitPrice    int           `json:"unitPrice"`
	SpecialPrice *SpecialOffer `json:"specialPrice"`
}

// SpecialOffer defines a multi-buy promotion.
type SpecialOffer struct {
	Quantity int `json:"quantity"`
	Price    int `json:"price"`
}
