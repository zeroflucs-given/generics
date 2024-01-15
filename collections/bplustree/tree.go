package bplustree

import (
	"fmt"
	"io"
	"sync"

	"github.com/zeroflucs-given/generics"
	"github.com/zeroflucs-given/generics/collections"
	"github.com/zeroflucs-given/generics/collections/ringbuffer"
)

const (
	MinTreeOrder = 2
)

// New creates a new instance of the B+ tree with the specified order.
func New[K generics.Comparable, V any](order int, preallocateSize int) (collections.TreeMap[K, V], error) {
	if order < MinTreeOrder {
		return nil, fmt.Errorf("invalid tree order %d: too low", order)
	} else if preallocateSize < 0 {
		return nil, fmt.Errorf("invalid pre-allocate size: %d too low", preallocateSize)
	}

	return &tree[K, V]{
		Order:                  order,
		preallocateSize:        preallocateSize,
		preallocatedKeySets:    ringbuffer.New[[]K](preallocateSize),
		preallocatedRecordSets: ringbuffer.New[[]record[V]](preallocateSize),
		preallocatedChildSets:  ringbuffer.New[[]*treeNode[K, V]](preallocateSize),
		preallocatedNodes:      ringbuffer.New[*treeNode[K, V]](preallocateSize),
	}, nil
}

type tree[K generics.Comparable, V any] struct {
	NodeCount              int64                                `json:"node_count"`   // Sequence number for allocating node
	RecordCount            collections.RecordID                 `json:"record_count"` // Record counter
	Order                  int                                  `json:"order"`        // Number of values in the tree
	Root                   *treeNode[K, V]                      `json:"root"`         // Root node
	lock                   sync.RWMutex                         `json:"-"`            // Lock to prevent concurrent modifies
	preallocateSize        int                                  `json:"-"`            // Pre-allocation/node pool sizes
	preallocatedKeySets    collections.Queue[[]K]               `json:"-"`            // Pre-allocation of key slices
	preallocatedRecordSets collections.Queue[[]record[V]]       `json:"-"`            // Pre-allocation of value slices
	preallocatedChildSets  collections.Queue[[]*treeNode[K, V]] `json:"-"`            // Pre-allocated child sets
	preallocatedNodes      collections.Queue[*treeNode[K, V]]   `json:"-"`            // Pre-allocated nodes
}

// Dump writes the tree out for diagnostic purposes to a file
func (t *tree[K, V]) Dump(f io.Writer) {
	_, _ = fmt.Fprintf(f, "========== TREE DUMP (Order: %v, NodesTotal: %v, Records: %v) ===========\n", t.Order, t.NodeCount, t.RecordCount)
	t.Root.Dump(f, 1)
}

// createNode creates a ne wnode in the tree
func (t *tree[K, V]) createNode(leaf bool) *treeNode[K, V] {
	t.NodeCount++
	nodeID := t.NodeCount

	result := &treeNode[K, V]{
		ID:   nodeID,
		Leaf: leaf,
		Keys: t.allocKeySet(),
	}

	// Pre-allocate appropriate child type
	if leaf {
		result.Records = t.allocRecordSet()
	} else {
		result.Children = t.allocChildSet()
	}

	return result
}

// allocKeySet allocates a key-set
func (t *tree[K, V]) allocKeySet() []K {
	if t.preallocateSize <= 1 {
		return make([]K, t.Order)
	}

	has, preAlloc := t.preallocatedKeySets.Pop()
	if has {
		return preAlloc
	}

	// Pre-allocate a large set and carve it up into blocks
	bigAlloc := make([]K, t.Order*t.preallocateSize)
	for i := 0; i < t.preallocateSize; i++ {
		firstIndex := i * t.Order
		err := t.preallocatedKeySets.Push(bigAlloc[firstIndex : firstIndex+t.Order])
		if err != nil {
			panic(err)
		}
	}

	_, first := t.preallocatedKeySets.Pop()
	return first
}

func (t *tree[K, V]) allocRecordSet() []record[V] {
	if t.preallocateSize <= 1 {
		return make([]record[V], t.Order)
	}

	has, preAlloc := t.preallocatedRecordSets.Pop()
	if has {
		return preAlloc
	}

	// Pre-allocate a large set and carve it up into blocks
	bigAlloc := make([]record[V], t.Order*t.preallocateSize)
	for i := 0; i < t.preallocateSize; i++ {
		firstIndex := i * t.Order
		subSlice := bigAlloc[firstIndex : firstIndex+t.Order]
		err := t.preallocatedRecordSets.Push(subSlice)
		if err != nil {
			panic(err)
		}
	}

	_, first := t.preallocatedRecordSets.Pop()
	return first
}

func (t *tree[K, V]) allocChildSet() []*treeNode[K, V] {
	if t.preallocateSize <= 1 {
		return make([]*treeNode[K, V], t.Order)
	}

	has, preAlloc := t.preallocatedChildSets.Pop()
	if has {
		return preAlloc
	}

	// Pre-allocate a large set and carve it up into blocks
	bigAlloc := make([]*treeNode[K, V], t.Order*t.preallocateSize)
	for i := 0; i < t.preallocateSize; i++ {
		firstIndex := i * t.Order
		err := t.preallocatedChildSets.Push(bigAlloc[firstIndex : firstIndex+t.Order])
		if err != nil {
			panic(err)
		}
	}

	_, first := t.preallocatedChildSets.Pop()
	return first
}

// findLeaf finds the insertion leaf node for a given key
func (t *tree[K, V]) findLeaf(k K) *treeNode[K, V] {
	current := t.Root
	// NB: We're using a linear scan here, same reason as treenode::getInsertIndex
	// Even when nodes have unreasonable sizes, the pay-off was non-existent versus
	// a linear scan (calculating pivots and branching costs more than iteration)
	// Whilst there is theoretically an inflection point, its just not worth it for
	// real-world use.

	// Keep recursing until we find a records/leaf node
	for current.Records == nil {
		// Assume we're going at the end
		targetIndex := current.Count - 1

		// Using a range here avoids repeated check of array bounds
		i := 1
		for _, currentKey := range current.Keys[i:current.Count] {
			if currentKey > k {
				targetIndex = i - 1
				break
			}
			i++
		}

		current = current.Children[targetIndex]
	}

	return current
}
