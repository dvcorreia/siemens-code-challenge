package app

import (
	"errors"
	"math/rand"
	"sync"
	"unicorn"
	"unicorn/pkg/queue"
	"unicorn/storage"
)

// Length of a Order ID.
// Exported so that it can be changed by developers.
var OrderIDLength = 16

var (
	ErrOrderCompleted = errors.New("order completed")
)

type order struct {
	ID unicorn.OrderID

	mu        sync.RWMutex
	amount    int // of unicorns to fullfil this order.
	produced  int // how many unicorn have been produced for this order.
	collected int // how many unicorns have been collected from this order.

	ready *queue.Queue[*unicorn.Unicorn] // unicorn ready for been collected.
}

// NewOrder creates a new unicorn production order.
func NewOrder(amount int) (*order, error) {
	id, err := randomID(OrderIDLength)
	if err != nil {
		return nil, err
	}

	return &order{
		ID:     unicorn.OrderID(id),
		amount: amount,
		ready:  queue.New[*unicorn.Unicorn](),
	}, nil
}

// Add unicorn to order.
func (o *order) Add(unicorn *unicorn.Unicorn) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.completed() {
		return ErrOrderCompleted
	}

	o.ready.Enqueue(unicorn)
	o.produced++
	return nil
}

// Completed indicates if the production for this order has been completed.
func (o *order) Completed() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.completed()
}

func (o *order) completed() bool {
	return o.amount == o.produced
}

// IsFulfilled returns the order fulfilled status and the number of pending unicorn.
func (o *order) IsFulfilled() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.amount == o.collected
}

// Pending returns the number of unicorns that are left to fulfill the order.
func (o *order) Pending() int {
	return o.amount - o.collected
}

// Collect available unicorns.
func (o *order) Collect() ([]*unicorn.Unicorn, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.ready.Empty() {
		return []*unicorn.Unicorn{}, nil
	}

	if o.completed() {
		return nil, ErrOrderCompleted
	}

	unicorns := o.ready.DequeueAll()
	o.collected += len(unicorns)

	return unicorns, nil
}

// CollectFromStorage fulfills the order with unicorns in storage.
func (o *order) CollectFromStorage(store storage.UnicornStorage) ([]*unicorn.Unicorn, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.completed() {
		return nil, ErrOrderCompleted
	}

	pending := o.amount - o.collected

	var toCollect int
	switch stored := store.InStorage(); {
	case stored == 0:
		return []*unicorn.Unicorn{}, nil
	case stored < pending:
		toCollect = stored
	case stored >= pending:
		toCollect = pending
	}

	o.produced += toCollect
	o.collected += toCollect
	return store.Collect(toCollect), nil
}

const charset = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ1234567890"

// randomID generates a random identifier n characters long.
func randomID(n int) (string, error) {
	b := make([]byte, n)

	for i := 0; i < len(b); i++ {
		b[i] = charset[rand.Intn(len(charset)-1)]
	}

	return string(b), nil
}
