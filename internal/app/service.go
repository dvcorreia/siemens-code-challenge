package app

import (
	"fmt"
	"sync"
	"unicorn"
	"unicorn/storage"
)

type service struct {
	mu sync.RWMutex

	store      storage.UnicornStorage
	production *productionLine

	orders map[unicorn.OrderID]*order
}

// New creates a new unicorn service app.
func New(production *productionLine, store storage.UnicornStorage) *service {
	return &service{
		store:      store,
		production: production,
		orders:     make(map[unicorn.OrderID]*order),
	}
}

// OrderUnicorns initiates a new unicorn production request.
// If no sufficient unicorn are available, it returns a request ID for consequent pooling.
func (s *service) OrderUnicorns(amount int) (unicorn.OrderID, error) {
	if amount <= 0 {
		return "", fmt.Errorf("invalid unicorn amount of %d", amount)
	}

	order, err := NewOrder(amount)
	if err != nil {
		return "", fmt.Errorf("generating new order: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.ID] = order

	if err := s.production.PlaceOrder(order); err != nil {
		return "", err
	}

	return order.ID, nil
}

// Pool returns the available ordered unicorns and how many are left to produce.
func (s *service) Pool(id unicorn.OrderID) ([]*unicorn.Unicorn, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[id]
	if !ok {
		return nil, ErrInvalidOrder
	}

	if order.IsFulfilled() {
		delete(s.orders, order.ID)
		return []*unicorn.Unicorn{}, nil
	}

	unicorns, err := order.Collect()
	if err != nil {
		return nil, err
	}

	if pending := order.Pending(); pending > 0 {
		uu, err := order.CollectFromStorage(s.store)
		if err != nil {
			return unicorns, err
		}
		unicorns = append(unicorns, uu...)
	}

	if order.IsFulfilled() {
		delete(s.orders, order.ID)
		return unicorns, nil
	}

	return unicorns, nil
}

// Validate checks if an ID has an orden in the process.
func (s *service) Validate(id unicorn.OrderID) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.orders[id]
	return ok
}

// PendingUnicorns checks how many unicorns are left to fulfill the production order.
func (s *service) PendingUnicorns(id unicorn.OrderID) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[id]
	if !ok {
		return 0, ErrInvalidOrder
	}

	if order.IsFulfilled() {
		return 0, nil
	}

	return order.Pending(), nil
}
