package ringbuffer

import (
	"fmt"

	collections "github.com/zeroflucs-given/generics/collections"
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
}

// Capacity of the buffer
func (b *RingBuffer[T]) Capacity() int {
	return b.capacity
}

// Count the number of records in the buffer
func (b *RingBuffer[T]) Count() int {
	head := b.head
	cursor := b.cursor

	if head < cursor {
		head = head + b.capacity + 1
	}

	return head - cursor
}

// Peek an item from the ring buffer
func (b *RingBuffer[T]) Peek() (bool, T) {
	// Buffers is empty
	if b.cursor == b.head {
		var def T
		return false, def
	}

	result := b.data[b.cursor]

	return true, result
}

// Pop an item from the ring buffer
func (b *RingBuffer[T]) Pop() (bool, T) {
	// Buffers is empty
	if b.cursor == b.head {
		var def T
		return false, def
	}

	result := b.data[b.cursor]
	b.cursor = b.cursor + 1
	if b.cursor == b.capacity+1 {
		b.cursor = 0 // Wrap
	}

	return true, result
}

// Push an item into the ring-buffer. Returns an error if we overflow
// the buffer
func (b *RingBuffer[T]) Push(item T) error {
	newHead := b.head + 1
	if newHead == b.capacity+1 {
		newHead = 0
	}

	if newHead == b.cursor {
		return fmt.Errorf("cursor wrapped at index %d: data may be lost: %w", newHead, collections.ErrBufferFull)
	}

	b.data[b.head] = item
	b.head = newHead

	return nil
}
