package app

import (
	"context"
	"errors"
	"sync"
	"time"
	"unicorn"
	"unicorn/factory"
	"unicorn/pkg/queue"
	"unicorn/storage"
)

var (
	ErrInvalidOrder   = errors.New("invalid production order")
	ErrInvalidStorage = errors.New("invalid unicorn storage")
)

type productionLine struct {
	mu sync.RWMutex

	rate time.Duration

	currentOrder *order
	orderQueue   *queue.Queue[*order]

	factory factory.Factory

	storage storage.UnicornStorage
}

// ProductionLineOption is function used to customize the production line.
type ProductionLineOption func(*productionLine) error

// WithStorage sets up a storage to save excedent production of unicorns.
func WithStorage(store storage.UnicornStorage) ProductionLineOption {
	return func(pl *productionLine) error {
		if store == nil {
			return ErrInvalidStorage
		}

		pl.storage = store
		return nil
	}
}

// NewProductionLine sets up a new unicorn production line.
func NewProductionLine(
	rate time.Duration,
	factory factory.Factory,
	options ...ProductionLineOption,
) (*productionLine, error) {
	queue := queue.New[*order]()

	pl := &productionLine{
		rate:       rate,
		orderQueue: queue,
		factory:    factory,
	}

	for _, opt := range options {
		if opt != nil {
			if err := opt(pl); err != nil {
				return nil, err
			}
		}
	}

	return pl, nil
}

// PlaceOrder adds an unicorn order to the production line.
func (pl *productionLine) PlaceOrder(order *order) error {
	if order == nil {
		return ErrInvalidOrder
	}

	pl.mu.Lock()
	defer pl.mu.Unlock()

	if pl.currentOrder == nil {
		pl.currentOrder = order
		return nil
	}

	pl.orderQueue.Enqueue(order)
	return nil
}

// ProductionRate return the unicorn production rate for the production line.
func (pl *productionLine) ProductionRate() time.Duration {
	pl.mu.RLock()
	defer pl.mu.RUnlock()
	return pl.rate
}

// Start production line.
func (pl *productionLine) Start(ctx context.Context) {
	for rate := pl.ProductionRate(); ; {
		select {
		case <-ctx.Done():
			return
		case <-time.After(rate):
			unicorn := pl.factory.NewUnicorn()
			pl.fulfill(unicorn)
		}
	}
}

// fulfills unicorn to the correct production order.
func (pl *productionLine) fulfill(unicorn *unicorn.Unicorn) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	order := pl.selectOrder()

	if order == nil {
		if pl.storage != nil {
			pl.storage.Store(unicorn)
		}
		return
	}

	order.Add(unicorn)
}

// selectOrder selects an order to fulfill.
func (pl *productionLine) selectOrder() *order {
	if pl.currentOrder == nil {
		return pl.nextOrder()
	}

	if pl.currentOrder.IsFulfilled() {
		return pl.nextOrder()
	}

	return pl.currentOrder
}

// select the next order to be produced.
func (pl *productionLine) nextOrder() *order {
	if pl.orderQueue.Empty() {
		pl.currentOrder = nil
		return pl.currentOrder
	}

	pl.currentOrder = pl.orderQueue.Dequeue()
	return pl.currentOrder
}
