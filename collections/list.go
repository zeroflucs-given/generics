package collections

const (
	// IndexNotFound is the index of an item that was not found
	IndexNotFound = -1
)

// List describes the behaviours of a generic list. Lists can have items inserted,
// removed or looked up by index.
type List[T comparable] interface {
	// Contains returns true if the specified item exists within the list.
	Contains(v T) bool

	// Insert adds an item to the list. An error is thrown if the insertion violates
	// a constraint/invariant of the list (such as a maximum capacity).
	Insert(v T) error

	// IndexOf gets the index of an item in the structure. If the item does not
	// exist in the collection then the return will be the IndexNotFound constant. If
	// multiple instances of the value exist, the first is returned, based on the
	// natural order of the structure.
	IndexOf(v T) int

	// Remove removes all instances of the specified item from the list.
	Remove(v T)

	// RemoveAt removes the item at the specified inde from the list. If the index
	// is beyond the bounds of the list, the operation will panic.
	RemoveAt(index int)

	// Value gets the value at a specified index in the structure. The boolean
	// value indicates if a value was returned or not. Values beyond the maximum
	// size of the list will return false and the default value of the type.
	Value(index int) (bool, T)
}
