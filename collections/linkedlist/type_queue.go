package linkedlist

import (
	collections "github.com/zeroflucs-given/generics/collections"
)

// Ensure our LinkedList implements the generic Queue[T] interface at compile time.
var _ collections.Queue[int] = (*LinkedList[int])(nil)

// Peek a value from the list
func (l *LinkedList[T]) Peek() (bool, T) {
	var found bool
	var result T

	l.lock.RLock()

	if l.head != nil {
		result = l.head.value
		found = true
	}

	l.lock.RUnlock()

	return found, result
}

// Pop a value from the list
func (l *LinkedList[T]) Pop() (bool, T) {
	var found bool
	var result T

	l.lock.Lock()

	if l.head != nil {
		result = l.head.value
		l.head = l.head.tail // Move backward
		found = true
	}

	l.lock.Unlock()

	return found, result
}

// Push a value into the linked list
func (l *LinkedList[T]) Push(value T) error {
	l.appendInternal(value)
	return nil
}
