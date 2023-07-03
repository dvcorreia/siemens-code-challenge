package app

import (
	"sync"
	"unicorn"
	"unicorn/factory"
	"unicorn/storage"
)

type ProductionLine <-chan *unicorn.Unicorn

type Order struct {
	Left           int
	ProductionLine ProductionLine
}

type service struct {
	mu sync.RWMutex

	factory factory.Factory
	storage storage.UnicornStorage

	orders map[unicorn.OrderID]Order
}

// New creates a new unicorn service app.
func New(factory factory.Factory, storage storage.UnicornStorage) *service {
	return &service{
		factory: factory,
		storage: storage,
		orders:  make(map[unicorn.OrderID]Order),
	}
}

// RequestUnicorns initiates a new unicorn production request.
// If no sufficient unicorn are available, it returns a request ID for consequent pooling.
func (s *service) OrderUnicorns(n int) (unicorn.OrderID, error) {
	return "", nil
}

// Pool returns the available ordered unicorns and how many are left to produce.
func (s *service) Pool(unicorn.OrderID) (int, []unicorn.Unicorn) {
	return 0, nil
}

// Validate checks if an ID has an orden in the process.
func (s *service) Validate(unicorn.OrderID) bool {
	return false
}
