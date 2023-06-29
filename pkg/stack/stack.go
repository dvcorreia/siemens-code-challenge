// Package stack provides an implementation of a LIFO stack built using a
// resizing array.
// Partial implementation of the https://github.com/zyedidia/generic stack.
package stack

// Stack implements a LIFO stack.
type Stack[T any] struct {
	entries []T
}

// New returns an empty stack.
func New[T any]() *Stack[T] {
	return &Stack[T]{
		entries: nil,
	}
}

// Push places 'value' at the top of the stack.
func (s *Stack[T]) Push(value T) {
	s.entries = append(s.entries, value)
}

// Pop removes the stack's top element and returns it. If the stack is empty it
// returns the zero value.
func (s *Stack[T]) Pop() (t T) {
	if len(s.entries) == 0 {
		return t
	}
	v := s.entries[len(s.entries)-1]
	s.entries = s.entries[:len(s.entries)-1]
	return v
}

// Size returns the number of elements in the stack.
func (s *Stack[T]) Size() int {
	return len(s.entries)
}
