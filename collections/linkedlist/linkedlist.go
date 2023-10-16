package linkedlist

import (
	"sync"

	"github.com/zeroflucs-given/generics/collections"
)

// New creates a new linked list.
func New[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{}
}

// LinkedList is our internal type for implementing the buffer pattern
type LinkedList[T comparable] struct {
	head *node[T]
	tail *node[T]
	lock sync.RWMutex
}

// Capacity of this linked list
func (l *LinkedList[T]) Capacity() int {
	return collections.CapacityInfinite
}

// Count of items in the list
func (l *LinkedList[T]) Count() int {
	count := 0

	l.lock.RLock()

	// Count backwards from the head to the tail
	current := l.head
	for current != nil {
		count++
		current = current.tail
	}

	l.lock.RUnlock()

	return count
}

func (l *LinkedList[T]) appendInternal(value T) {
	l.lock.Lock()

	// Build our new node and join to the tail
	oldTail := l.tail
	newTail := &node[T]{
		head:  oldTail,
		value: value,
	}
	if oldTail != nil {
		oldTail.tail = newTail
	}
	l.tail = newTail

	// If we're the first value
	if l.head == nil {
		l.head = newTail
	}

	l.lock.Unlock()
}

// node is a node in the linked list
type node[T any] struct {
	head  *node[T]
	tail  *node[T]
	value T
}
