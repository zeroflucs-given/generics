package linkedlist

import (
	"fmt"

	"github.com/zeroflucs-given/generics/collections"
)

// Ensure our LinkedList implements the generic List[T] interface at compile time.
var _ collections.List[int] = (*LinkedList[int])(nil)

// Contains returns true if the specified item exists in the list
func (l *LinkedList[T]) Contains(v T) bool {
	return l.IndexOf(v) != collections.IndexNotFound
}

// Insert appends an item to the list
func (l *LinkedList[T]) Insert(v T) error {
	l.appendInternal(v)
	return nil
}

// Remove removes all instances of the specified value from the list
func (l *LinkedList[T]) Remove(v T) {
	l.lock.Lock()

	var previousTail *node[T]
	current := l.tail
	for current != nil {
		// Cut this item out of the list?
		if current.value == v {
			// If we have an item before us, then link its next
			// node to *our* next node
			if previousTail != nil {
				previousTail.head = current.head
			}

			// If there is an item after us, then its previous is
			// now our previous
			if current.head != nil {
				current.head.tail = previousTail
			}

			if current == l.head {
				l.head = current.tail
			}
			if current == l.tail {
				l.tail = current.head
			}
		}

		// Move forward in the list
		previousTail = current
		current = current.head
	}

	l.lock.Unlock()
}

// RemoveAt removes the item with the specified index from the list
func (l *LinkedList[T]) RemoveAt(index int) {
	l.lock.Lock()

	var previousTail *node[T]
	var currentIndex int

	current := l.tail
	for current != nil {
		// Cut this item out of the list?
		if currentIndex == index {
			// If we have an item before us, then link its next
			// node to *our* next node
			if previousTail != nil {
				previousTail.head = current.head
			}

			// If there is an item after us, then its previous is
			// now our previous
			if current.head != nil {
				current.head.tail = previousTail
			}

			if current == l.head {
				l.head = current.tail
			}

			if current == l.tail {
				l.tail = current.head
			}
		}

		// Done
		if currentIndex >= index {
			break
		}

		// Move forward in the list
		previousTail = current
		current = current.head
		if current != nil {
			currentIndex++
		}
	}

	l.lock.Unlock()

	if currentIndex < index {
		panic(fmt.Sprintf("The index %d is beyond the bounds of the LinkedList", index))
	}
}

// IndexOf gets the index of an item in the structure
func (l *LinkedList[T]) IndexOf(v T) int {
	result := collections.IndexNotFound
	var currentIndex int

	l.lock.RLock()
	current := l.tail
	for current != nil {
		if current.value == v {
			result = currentIndex
			break
		}
		current = current.head
		currentIndex++
	}
	l.lock.RUnlock()

	return result
}

// Value gets a value by its index in the list
func (l *LinkedList[T]) Value(index int) (bool, T) {
	var found bool
	var result T
	var currentIndex int

	l.lock.RLock()
	current := l.tail
	for current != nil {
		if currentIndex == index {
			found = true
			result = current.value
			break
		}
		current = current.head
		currentIndex++
	}
	l.lock.RUnlock()

	return found, result
}
