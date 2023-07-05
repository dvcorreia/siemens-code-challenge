package app

import (
	"math/rand"
	"sync"
	"unicorn"
	"unicorn/pkg/queue"
)

// Length of a Order ID.
// Exported so that it can be changed by developers.
var OrderIDLength = 16

type order struct {
	ID unicorn.OrderID

	mu       sync.RWMutex
	amount   int // of unicorns to fullfil this order.
	produced int // how many unicorn have been produced for this order.
	sent     int // how many unicorns have been shipped to clients.

	ready *queue.Queue[*unicorn.Unicorn] // unicorn ready for been collected.
}

// NewOrder creates a new unicorn production order.
func NewOrder(amount uint) *order {
	id := randomID(OrderIDLength)

	return &order{
		ID:     unicorn.OrderID(id),
		amount: int(amount),
		ready:  queue.New[*unicorn.Unicorn](),
	}
}

// Collect available unicorns.
func (o *order) Collect() []*unicorn.Unicorn {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.ready.Empty() {
		return []*unicorn.Unicorn{}
	}

	unicorns := o.ready.DequeueAll()
	o.sent += len(unicorns)

	return unicorns
}

// Add unicorn to order. It returns ok.
func (o *order) Add(unicorn *unicorn.Unicorn) bool {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.amount == o.produced {
		return false
	}

	o.ready.Enqueue(unicorn)
	o.produced++
	return true
}

// ProductionHasCompleted indicates if the production for this order has been completed.
func (o *order) ProductionHasCompleted() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.amount == o.produced
}

// IsFulfilled indicates if the order has been completed.
func (o *order) IsFulfilled() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.amount == o.sent
}

// PendingProduction returns the number of unicorns that are left to fulfill the order.
func (o *order) PendingProduction() int {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.amount - o.produced
}

const charset = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ1234567890"

// randomID generates a random identifier n characters long.
func randomID(n int) string {
	b := make([]byte, n)

	for i := 0; i < len(b); i++ {
		b[i] = charset[rand.Intn(len(charset)-1)]
	}

	return string(b)
}
