package stack

import collections "github.com/zeroflucs-given/generics/collections"

// Ensure we meet the Queue[T] interface at compile time
var _ collections.Queue[int] = (*Stack[int])(nil)

// Push a value into the stack
func (s *Stack[T]) Push(value T) error {
	s.lock.Lock()

	if len(s.data) == s.head {
		s.lock.Unlock()

		return collections.ErrBufferFull
	}

	s.data[s.head] = value
	s.head = s.head + 1

	s.lock.Unlock()
	return nil
}

// Peek the item that would be returned from pop. Note that for LIFO situations the
// list consumer must mutex, otherwise the next Pop may yield a different element
func (s *Stack[T]) Peek() (bool, T) {
	found := false
	var v T

	s.lock.RLock()

	if s.head > 0 {
		dataIndex := s.head - 1
		v = s.data[dataIndex]
		found = true
	}

	s.lock.RUnlock()

	return found, v
}

// Pop a value from the stack
func (s *Stack[T]) Pop() (bool, T) {
	found := false
	var v T
	var blank T

	s.lock.Lock()

	if s.head > 0 {
		dataIndex := s.head - 1
		v = s.data[dataIndex]
		found = true
		s.data[dataIndex] = blank
		s.head = s.head - 1
	}

	s.lock.Unlock()

	return found, v
}
