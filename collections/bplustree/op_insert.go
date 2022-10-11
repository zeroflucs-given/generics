package bplustree

import collections "github.com/zeroflucs-given/generics/collections"

// Insert a value into the tree
func (t *tree[K, V]) Insert(key K, value V) collections.RecordID {
	t.lock.Lock()

	// Allocate a record ID
	t.RecordCount++
	recordID := t.RecordCount
	record := record[V]{
		RecordID: recordID,
		Value:    value,
	}

	// Case: Empty tree
	if t.Root == nil {
		root := t.createNode(true)
		root.insertRecord(key, record)
		t.Root = root
		t.lock.Unlock()
		return recordID
	}

	targetLeaf := t.findLeaf(key)
	if targetLeaf.Count == t.Order {
		targetLeaf = t.split(targetLeaf, key, 0)
	}

	// Write to the records list
	targetLeaf.insertRecord(key, record)

	t.lock.Unlock()
	return recordID
}

// split a node to make room for friends
func (t *tree[K, V]) split(existingNode *treeNode[K, V], keyToAccomodate K, depth int) *treeNode[K, V] {
	// Track some state before things get exciting
	existingParent := existingNode.Parent

	// Split the children/data of the subject node
	newSibling := t.createNode(existingNode.Leaf)
	splitPoint := existingNode.Count / 2

	targetIndex := 0
	for i := splitPoint; i < existingNode.Count; i++ {
		newSibling.Keys[targetIndex] = existingNode.Keys[i]
		if existingNode.Leaf {
			newSibling.Records[targetIndex] = existingNode.Records[i]
		} else {
			child := existingNode.Children[i]
			newSibling.Children[targetIndex] = child
			existingNode.Children[i] = nil
			child.Parent = newSibling
		}
		targetIndex++
	}

	existingNode.Count = splitPoint
	newSibling.Count = t.Order - splitPoint
	newSiblingFirstKey := newSibling.Keys[0]

	if existingParent == nil {
		// Case 1 - We're splitting the root
		newRoot := t.createNode(false)
		newRoot.insertChild(existingNode.Keys[0], existingNode)
		newRoot.insertChild(newSiblingFirstKey, newSibling)
		t.Root = newRoot
	} else if existingParent.Count < t.Order {
		// Case 2 - We're inserting into a parent that has space
		existingParent.insertChild(newSiblingFirstKey, newSibling)
	} else {
		// Case 3 - Split recursively
		newParent := t.split(existingNode.Parent, newSiblingFirstKey, depth+1)
		newParent.insertChild(newSiblingFirstKey, newSibling)
	}

	// Now determine which of the two nodes we call home
	if keyToAccomodate < newSiblingFirstKey {
		return existingNode
	}

	return newSibling
}
