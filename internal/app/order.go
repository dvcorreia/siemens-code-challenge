package app

import (
	"crypto/rand"
	"errors"
	"sync"
	"unicorn"
	"unicorn/pkg/queue"
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
	return o.amount <= o.produced
}

// IsFulfilled returns the order fulfilled status and the number of pending unicorn.
func (o *order) IsFulfilled() (int, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	pending := o.amount - o.collected

	return pending, pending <= 0
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

// randomID generates a random identifier n characters long.
func randomID(n int) (string, error) {
	c := 16
	b := make([]byte, c)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return string(b), nil
}
