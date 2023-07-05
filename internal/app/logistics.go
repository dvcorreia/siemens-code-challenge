package app

import (
	"sync"
	"unicorn"
	"unicorn/pkg/queue"
	"unicorn/storage"
)

type logisticsCenter struct {
	mu sync.RWMutex

	current *order                 // current order to fulfill. must be always non nil.
	queue   *queue.Queue[*order]   // queue of order to fulfill
	store   storage.UnicornStorage // unicorn store for excedent production
}

func NewLogisticsCenter(store storage.UnicornStorage) *logisticsCenter {
	queue := queue.New[*order]()

	return &logisticsCenter{
		current: NewOrder(0), // select a fake zero amount order. This facilitates the logic.
		queue:   queue,
		store:   store,
	}
}

func (lc *logisticsCenter) AddOrder(order *order) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if lc.store.InStorage() > 0 {
		pending := order.amount - order.produced

		unicorns := lc.store.Collect(pending)

		for _, u := range unicorns {
			if !order.Add(u) {
				lc.store.Store(u)
			}
		}
	}

	lc.queue.Enqueue(order)
}

func (lc *logisticsCenter) HandleUnicorn(unicorn *unicorn.Unicorn) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.updateCurrentOrder()

	if !lc.current.Add(unicorn) {
		lc.store.Store(unicorn)
	}
}

func (lc *logisticsCenter) updateCurrentOrder() {
	// if the production has not ended, keep it in production.
	if !lc.current.ProductionHasCompleted() {
		return
	}

	// in case that no other orders are available, keep the current one.
	if lc.queue.Empty() {
		return
	}

	lc.current = lc.queue.Dequeue()
}
