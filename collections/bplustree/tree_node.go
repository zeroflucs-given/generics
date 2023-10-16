package bplustree

import (
	"fmt"
	"io"
	"strings"

	"github.com/zeroflucs-given/generics"
)

// treeNode is a node within the tree
type treeNode[K generics.Comparable, V any] struct {
	ID         int64  `json:"node_id"`    // Unique node ID
	Leaf       bool   `json:"leaf"`       // Leaf/Data node?
	Keys       []K    `json:"key"`        // Keys
	Count      int    `json:"count"`      // Number of children/data records
	Annotation string `json:"annotation"` // Annotation/Informational tag

	// Genealogy
	Parent          *treeNode[K, V] `json:"-"` // Parent
	PreviousSibling *treeNode[K, V] `json:"-"` // Previous sibling
	NextSibling     *treeNode[K, V] `json:"-"` // Next sibling

	// Storage
	Children []*treeNode[K, V] `json:"children"` // Child nodes
	Records  []record[V]       `json:"records"`  // Records
}

func (tn *treeNode[K, V]) Dump(f io.Writer, depth int) {
	prefix := strings.Repeat("\t|", depth) + "-"
	if tn == nil {
		fmt.Fprintf(f, "%v(Empty)\n", prefix)
		return
	}

	fmt.Fprintf(f, "%vNODE(%v) [Leaf=%v, Count=%v] %q\n", prefix, tn.ID, tn.Leaf, tn.Count, tn.Annotation)
	fmt.Fprintf(f, "%v  Prev: %v | Parent: %v | Next: %v\n", prefix, tn.PreviousSibling.NodeID(), tn.Parent.NodeID(), tn.NextSibling.NodeID())

	for i := 0; i < tn.Count; i++ {
		if tn.Leaf {
			fmt.Fprintf(f, "%v- RECORD %d [Key: %v, Record #: %d]\n", prefix, i, tn.Keys[i], tn.Records[i].RecordID)
		} else {
			fmt.Fprintf(f, "%v- CHILD #%d [Key: %v] NODE(%v)\n", prefix, i, tn.Keys[i], tn.Children[i].ID)
			tn.Children[i].Dump(f, depth+1)
		}
	}

}

func (tn *treeNode[K, V]) NodeID() string {
	if tn == nil {
		return "(none)"
	}
	return fmt.Sprintf("%v", tn.ID)
}

// indexOf gets the index of the specified child in the slice
func (tn *treeNode[K, V]) indexOf(subject *treeNode[K, V]) int {
	if tn != nil && subject != nil {
		for i, current := range tn.Children[0:tn.Count] {
			if current == subject {
				return i
			}
		}
	}

	return -1
}

// insertChild inserts a child into the node.
func (tn *treeNode[K, V]) insertChild(key K, child *treeNode[K, V]) {
	var previousChildSibling, nextChildSibling *treeNode[K, V]
	if tn.Count > 0 {
		previousChildSibling = tn.Children[0].PreviousSibling
		nextChildSibling = tn.Children[tn.Count-1].NextSibling
	}

	targetIndex := tn.getInsertIndex(key)

	// Move all later values down from the read
	for lastIndex := tn.Count - 1; lastIndex >= targetIndex; lastIndex-- {
		tn.Keys[lastIndex+1] = tn.Keys[lastIndex]
		tn.Children[lastIndex+1] = tn.Children[lastIndex]
	}

	// Add the value the index
	tn.Keys[targetIndex] = key
	tn.Children[targetIndex] = child
	tn.Count++

	child.Parent = tn

	// Maintain positions
	for i, current := range tn.Children[0:tn.Count] {

		// If we're the first node then link to the original predecessor
		if i == 0 {
			current.PreviousSibling = previousChildSibling
			if previousChildSibling != nil {
				previousChildSibling.NextSibling = current
			}
		}

		// If we're the last node, then link to the original successor
		if i == tn.Count-1 {
			current.NextSibling = nextChildSibling
			if nextChildSibling != nil {
				nextChildSibling.PreviousSibling = current
			}
		}

		// If we're internal, link previous siblings to prior nodes
		if i > 0 {
			current.PreviousSibling = tn.Children[i-1]
		}
		if i < tn.Count-1 {
			current.NextSibling = tn.Children[i+1]
		}
	}

	if targetIndex == 0 {
		parentIndex := tn.Parent.indexOf(tn)
		if parentIndex >= 0 {
			tn.Parent.Keys[parentIndex] = key
		}
	}
}

// insertRecord inserts a record into the node.
func (tn *treeNode[K, V]) insertRecord(key K, record record[V]) {
	targetIndex := tn.getInsertIndex(key)

	// Move all later values down from the read
	for lastIndex := tn.Count - 1; lastIndex >= targetIndex; lastIndex-- {
		tn.Keys[lastIndex+1] = tn.Keys[lastIndex]
		tn.Records[lastIndex+1] = tn.Records[lastIndex]
	}

	tn.Keys[targetIndex] = key
	tn.Records[targetIndex] = record
	tn.Count++

	// If we're inserting at the lead slot, update our
	// parents key for us
	if targetIndex == 0 {
		tn.updateParentReference()
	}
}

func (tn *treeNode[K, V]) updateParentReference() {
	if tn.Parent == nil {
		return
	}

	parentIndex := tn.Parent.indexOf(tn)
	leadKey := tn.Keys[0]

	tn.Parent.Keys[parentIndex] = leadKey

	if parentIndex == 0 {
		tn.Parent.updateParentReference()
	}
}

// getInsertIndex gets the insertion index for value K
func (tn *treeNode[K, V]) getInsertIndex(k K) int {
	// Sequential scan
	var targetIndex int
	for _, nodeKey := range tn.Keys[0:tn.Count] {
		if nodeKey > k {
			break
		}
		targetIndex++
	}

	return targetIndex

	/*
		NB: There was no improvement using a binary search here.
		The code for this is below, however in practice even with 65k
		entries per node, the extra calculations required negated any
		improvement.

			keys := tn.Keys
			low := 0
			high := tn.Count - 1

			for low <= high {
				median := (low + high) / 2
				mk := keys[median]
				if mk < k {
					low = median + 1 // Second half
				} else {
					high = median - 1 // First half
				}
			}
			return low
	*/
}
