package pricing

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/TheFodfather/checkoutapi/domain"
)

// Service provides access to pricing rules
type Service struct {
	pricingFile string
	rules       map[string]domain.PricingRule
	lastModTime time.Time
	sync.RWMutex
}

// New creates a new pricing service and loads pricing data
func New(pricingFilePath string) (*Service, error) {
	s := &Service{
		pricingFile: pricingFilePath,
		rules:       make(map[string]domain.PricingRule),
	}

	if err := s.loadPricingRules(); err != nil {
		return nil, fmt.Errorf("initial pricing load failed: %w", err)
	}

	go s.watchPricingChanges()

	return s, nil
}

// GetRules returns a copy of the current pricing rules.
func (s *Service) GetRules() map[string]domain.PricingRule {
	s.RLock()
	defer s.RUnlock()

	rulesCopy := make(map[string]domain.PricingRule)

	for k, v := range s.rules {
		rulesCopy[k] = v
	}

	return rulesCopy
}

func (s *Service) loadPricingRules() error {
	file, err := os.ReadFile(s.pricingFile)
	if err != nil {
		return err
	}

	var newRules map[string]domain.PricingRule
	if err := json.Unmarshal(file, &newRules); err != nil {
		return fmt.Errorf("failed to parse pricing json: %w", err)
	}

	fileInfo, err := os.Stat(s.pricingFile)
	if err != nil {
		return err
	}

	s.Lock()
	s.rules = newRules
	s.lastModTime = fileInfo.ModTime()
	s.Unlock()

	log.Println("✅ Successfully loaded new pricing rules.")

	return nil
}

func (s *Service) watchPricingChanges() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fileInfo, err := os.Stat(s.pricingFile)
		if err != nil {
			log.Printf("⚠️ Could not stat pricing file, keeping old rules: %v", err)
			continue
		}

		s.RLock()
		lastMod := s.lastModTime
		s.RUnlock()

		if fileInfo.ModTime().After(lastMod) {
			log.Println("🔄 Change detected in pricing.json, attempting to reload...")
			if err := s.loadPricingRules(); err != nil {
				log.Printf("❌ Error reloading pricing rules: %v", err)
			}
		}
	}
}
