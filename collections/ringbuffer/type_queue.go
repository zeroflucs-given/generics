package ringbuffer

import (
	"fmt"

	collections "github.com/zeroflucs-given/generics/collections"
)

// Ensure we meet the Queue[T] interface at compile time
var _ collections.Queue[int] = (*RingBuffer[int])(nil)

// Peek an item from the ring buffer
func (b *RingBuffer[T]) Peek() (bool, T) {
	b.lock.RLock()

	// Buffers is empty
	if b.cursor == b.head {
		var def T
		b.lock.RUnlock()
		return false, def
	}

	result := b.data[b.cursor]

	b.lock.RUnlock()

	return true, result
}

// Pop an item from the ring buffer
func (b *RingBuffer[T]) Pop() (bool, T) {
	b.lock.Lock()

	// Buffers is empty
	if b.cursor == b.head {
		var def T
		b.lock.Unlock()
		return false, def
	}

	result := b.data[b.cursor]
	b.cursor = b.cursor + 1
	if b.cursor == b.capacity+1 {
		b.cursor = 0 // Wrap
	}

	b.lock.Unlock()

	return true, result
}

// Push an item into the ring-buffer. Returns an error if we overflow
// the buffer
func (b *RingBuffer[T]) Push(item T) error {
	b.lock.Lock()

	newHead := b.head + 1
	if newHead == b.capacity+1 {
		newHead = 0
	}

	if newHead == b.cursor {
		b.lock.Unlock()
		return fmt.Errorf("cursor wrapped at index %d: data may be lost: %w", newHead, collections.ErrBufferFull)
	}

	b.data[b.head] = item
	b.head = newHead

	b.lock.Unlock()

	return nil
}
