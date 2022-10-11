package bplustree

import (
	"fmt"
	"os"
)

// CheckConsistency checks the consistency of the data-structure and ensures that there are
// no obvious problems with it.
func (t *tree[K, V]) CheckConsistency() {
	// Find the leftmost edge of the tree
	var defKey K
	node := t.findLeaf(defKey)
	for node.PreviousSibling != nil {
		node = node.PreviousSibling
	}

	levels := []*treeNode[K, V]{}
	for node != nil {
		levels = append(levels, node)
		node = node.Parent
	}
	for i, j := 0, len(levels)-1; i < j; i, j = i+1, j-1 {
		levels[i], levels[j] = levels[j], levels[i]
	}

	for level, leftMost := range levels {
		fmt.Printf("Consistency checking level %d: Leftmost NODE(%v)\n", level, leftMost.NodeID())
		t.checkLevelConsistency(leftMost)
	}

	fmt.Println("--- TREE CONSISTENT ---")
}

// checkLevelConsistency assumes that the node input is the first/leftmost edge of a tree level
// and walks along the tree, ensuring values are correctly managed.
func (t *tree[K, V]) checkLevelConsistency(node *treeNode[K, V]) {
	if node == nil {
		return
	}

	// Check all keys ordered
	currentKey := node.Keys[0]
	current := node
	orderOK := true
orderCheck:
	for current != nil {
		for i := 0; i < current.Count; i++ {
			visitKey := current.Keys[i]
			if visitKey < currentKey {
				current.Annotation = "FAIL"
				fmt.Printf("         NODE(%d) had a rewind. Key rewound from %v to %v!\n", current.ID, currentKey, visitKey)
				orderOK = false
				break orderCheck
			}
		}
		current = current.NextSibling
	}

	// Check prev/next links
	current = node
	sequenceOK := true
	for current != nil {
		// Our predecessor should link together
		if current.PreviousSibling != nil {
			if current.PreviousSibling.NextSibling != current {
				current.Annotation = "FAIL"
				fmt.Printf("         NODE(%d) Our previous sibling (%d) next does not link next to us. Points to (%v)!\n", current.ID, current.PreviousSibling.ID, current.PreviousSibling.NextSibling.NodeID())
				sequenceOK = false
			}

			if current.PreviousSibling.Parent != nil {
				childIndex := current.PreviousSibling.Parent.indexOf(current.PreviousSibling)
				if childIndex < 0 {
					current.PreviousSibling.Annotation = "ORPHAN"
					fmt.Printf("         NODE(%d) Our previous sibling (%d) is not registered to its parent %v\n", current.ID, current.PreviousSibling.ID, current.PreviousSibling.Parent.NodeID())
					sequenceOK = false
				}
			}
		}

		if current.NextSibling != nil {
			if current.NextSibling.PreviousSibling != current {
				current.Annotation = "FAIL"
				fmt.Printf("         NODE(%d) Our next sibling (%d) previous does not link next to us. Points to (%v)!\n", current.ID, current.NextSibling.ID, current.NextSibling.PreviousSibling.NodeID())
				sequenceOK = false
			}
		}

		current = current.NextSibling
	}

	if !orderOK || !sequenceOK {
		fmt.Println("==================== TREE DUMP =============== ")
		t.Dump(os.Stderr)
		panic("Inconsistent state")
	}
}
