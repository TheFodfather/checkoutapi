package repository

import (
	"fmt"
	"sync"

	"github.com/TheFodfather/checkoutapi/domain"
)

// SessionRepository defines the interface for storing checkout sessions.
type SessionRepository interface {
	Get(id string) (domain.ICheckout, error)
	Save(co domain.ICheckout) error
}

type InMemoryRepository struct {
	sessions map[string]domain.ICheckout
	sync.RWMutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		sessions: make(map[string]domain.ICheckout),
	}
}

func (m *InMemoryRepository) Get(id string) (domain.ICheckout, error) {
	m.RLock()
	defer m.RUnlock()
	co, ok := m.sessions[id]
	if !ok {
		return nil, fmt.Errorf("session with id '%s' not found", id)
	}
	return co, nil
}

func (m *InMemoryRepository) Save(co domain.ICheckout) error {
	m.Lock()
	defer m.Unlock()
	m.sessions[co.GetID()] = co
	return nil
}
