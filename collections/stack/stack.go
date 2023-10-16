package stack

import (
	"sync"
)

// NewStack creates a new instance of a stack with an initial capacity
func NewStack[T any](capacity int) *Stack[T] {
	return &Stack[T]{
		data: make([]T, capacity),
	}
}

// Stack is our type that implements a stack of data items
type Stack[T any] struct {
	data []T
	head int
	lock sync.RWMutex
}

// Count of the data inside the stack
func (s *Stack[T]) Count() int {
	s.lock.RLock()
	count := s.head
	s.lock.RUnlock()
	return count
}

// Capacity is the capacity of this stack. Stacks can grow until memory
// is depleted today.
func (s *Stack[T]) Capacity() int {
	return len(s.data)
}
