package ringbuffer

import (
	"sync"
)

// New is a fixed-size ring/circle buffer of values.
func New[T any](size int) *RingBuffer[T] {
	return &RingBuffer[T]{
		capacity: size,
		data:     make([]T, size+1), // We actually keep+1
	}
}

type RingBuffer[T any] struct {
	cursor   int
	capacity int
	head     int
	data     []T
	lock     sync.RWMutex
}

// Capacity of the buffer
func (b *RingBuffer[T]) Capacity() int {
	return b.capacity
}

// Count the number of records in the buffer
func (b *RingBuffer[T]) Count() int {
	b.lock.RLock()
	head := b.head
	cursor := b.cursor
	b.lock.RUnlock()

	if head < cursor {
		head = head + b.capacity + 1
	}

	return head - cursor
}
