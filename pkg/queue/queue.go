// Package queue provides an implementation of a First In First Out (FIFO) queue.
// The FIFO queue is implemented using the container/list standard library doubly-linked list
// and is based of the https://github.com/zyedidia/generic queue.
package queue

import "container/list"

// Queue is a simple First In First Out (FIFO) queue.
type Queue[T any] struct {
	list *list.List
}

// New returns an empty First In First Out (FIFO) queue.
func New[T any]() *Queue[T] {
	return &Queue[T]{
		list: list.New(),
	}
}

// Len returns the number of items currently in the queue.
func (q *Queue[T]) Len() int {
	return q.list.Len()
}

// Enqueue inserts 'value' to the end of the queue.
func (q *Queue[T]) Enqueue(value T) {
	q.list.PushBack(value)
}

// Dequeue removes and returns the item at the front of the queue.
//
// A panic occurs if the queue is Empty.
func (q *Queue[T]) Dequeue() T {
	value, ok := q.TryDequeue()
	if !ok {
		panic("queue: tried to dequeue from an empty queue")
	}
	return value
}

// TryDequeue tries to remove and return the item at the front of the queue.
//
// If the queue is empty, then false is returned as the second return value.
func (q *Queue[T]) TryDequeue() (T, bool) {
	if q.Empty() {
		var zero T
		return zero, false
	}

	value := q.list.Remove(q.list.Front()).(T)
	return value, true
}

// DequeueAll removes and returns all the items in the queue.
func (q *Queue[T]) DequeueAll() []T {
	slice := make([]T, q.list.Len())
	for i := 0; i < len(slice); i++ {
		slice[i] = q.Dequeue()
	}
	return slice
}

// Empty returns true if the queue is empty.
func (q *Queue[T]) Empty() bool {
	return q.list.Len() == 0
}
