package collections

// Queue is an interface that describes any LIFO/FIFO queue of values
type Queue[T any] interface {
	// Count the number of records in the buffer
	Count() int

	// Capacity of the buffer.
	Capacity() int

	// Peek the item that would be returned from pop. Note that for LIFO situations the
	// list consumer must mutex, otherwise the next Pop may yield a different element
	Peek() (bool, T)

	// Pop an item from the buffer. First argument indicates if the value was found
	// or not, allowing for T to be a value-type (i.e. int) where it may be hard to
	// distinguish empty based on the default (0 vs not there)
	Pop() (bool, T)

	// Push a new item into the buffer
	Push(t T) error
}
