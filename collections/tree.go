package collections

import (
	"github.com/zeroflucs-given/generics"
)

// RecordID is a unique identifier for a record in a tree. This allows us to
// differentiate tree records.
type RecordID int64

// TreeMap is our interface for a key-value map stored in a seekable tree format.
type TreeMap[K generics.Comparable, V any] interface {
	// Insert a value into the tree.
	Insert(key K, value V) RecordID

	// Scan records
	Scan() chan generics.KeyValuePair[K, V]

	// Count records
	Count() int
}
