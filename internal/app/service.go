package app

import (
	"errors"
	"fmt"
	"sync"
	"unicorn"
)

var ErrInvalidOrder = errors.New("invalid production order")

type service struct {
	mu sync.RWMutex

	logistics *logisticsCenter

	// to keep track of pending orders
	orders map[unicorn.OrderID]*order
}

// New creates a new unicorn service app.
func New(center *logisticsCenter) *service {
	return &service{
		logistics: center,
		orders:    make(map[unicorn.OrderID]*order),
	}
}

var _ unicorn.Service = (*service)(nil)

// OrderUnicorns initiates a new unicorn production request.
// If no sufficient unicorn are available, it returns a request ID for consequent pooling.
func (s *service) OrderUnicorns(amount int) (unicorn.OrderID, error) {
	if amount <= 0 {
		return "", fmt.Errorf("invalid unicorn amount of %d", amount)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	order := NewOrder(uint(amount))

	s.orders[order.ID] = order
	s.logistics.AddOrder(order)

	return order.ID, nil
}

// Pool returns the available ordered unicorns and how many are left to produce.
func (s *service) Pool(id unicorn.OrderID) ([]*unicorn.Unicorn, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[id]
	if !ok {
		return nil, 0, ErrInvalidOrder
	}

	unicorns := order.Collect()
	pending := order.PendingProduction()

	if order.IsFulfilled() {
		delete(s.orders, order.ID)
	}

	return unicorns, pending, nil
}

// Validate checks if an ID has an orden in the process.
func (s *service) Validate(id unicorn.OrderID) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.orders[id]
	return ok
}
