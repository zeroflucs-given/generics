package collections

import "errors"

const (
	// CapacityInfinite is a way of indicating the buffer has infinite
	// capacity, bound only by system resources.
	CapacityInfinite = -1
)

// ErrBufferFull indicates a buffer cannot be written to.
var ErrBufferFull = errors.New("the buffer is full and cannot take more data")
