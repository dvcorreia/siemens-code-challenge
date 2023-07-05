package app

import (
	"fmt"
	"sync"
	"unicorn"
)

type service struct {
	mu sync.RWMutex

	production *productionLine
	orders     map[unicorn.OrderID]*order
}

// New creates a new unicorn service app.
func New(production *productionLine) *service {
	return &service{
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

	return order.ID, nil
}

// Pool returns the available ordered unicorns and how many are left to produce.
func (s *service) Pool(id unicorn.OrderID) (int, []unicorn.Unicorn) {
	// grab from the order

	// del the order if its done

	return 0, nil
}

// Validate checks if an ID has an orden in the process.
func (s *service) Validate(id unicorn.OrderID) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.orders[id]
	return ok
}
