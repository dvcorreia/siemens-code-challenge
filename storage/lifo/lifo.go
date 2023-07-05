package lifo

import (
	"sync"
	"unicorn"
	"unicorn/pkg/stack"

	unicornstorage "unicorn/storage"
)

type storage struct {
	mu    sync.RWMutex
	stack *stack.Stack[*unicorn.Unicorn]
}

// New creates a unicorn LIFO store.
func New() *storage {
	return &storage{
		stack: stack.New[*unicorn.Unicorn](),
	}
}

var _ unicornstorage.UnicornStorage = (*storage)(nil)

// Store places a unicorn in storage.
func (s *storage) Store(unicorn *unicorn.Unicorn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stack.Push(unicorn)
}

// InStorage returns the number o unicorns in storage.
func (s *storage) InStorage() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stack.Size()
}

// Collect will do a best effort of collecting a number of unicorns from storage.
// If there are not enough unicorns in storage, it will return any it can provide.
func (s *storage) Collect(n int) []*unicorn.Unicorn {
	s.mu.Lock()
	defer s.mu.Unlock()

	var unicorns []*unicorn.Unicorn

	switch l := s.stack.Size(); {
	case l == 0:
		unicorns = []*unicorn.Unicorn{}
	case l <= n:
		unicorns = make([]*unicorn.Unicorn, l)
		for i := 0; i < l; i++ {
			unicorns[i] = s.stack.Pop()
		}
	case l > n:
		unicorns = make([]*unicorn.Unicorn, n)
		for i := 0; i < n; i++ {
			unicorns[i] = s.stack.Pop()
		}
	}

	return unicorns
}
