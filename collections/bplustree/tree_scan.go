package bplustree

import (
	"github.com/zeroflucs-given/generics"
)

func (t *tree[K, V]) Count() int {
	counter := 0
	for range t.Scan() {
		counter++
	}
	return counter
}

func (t *tree[K, V]) Scan() chan generics.KeyValuePair[K, V] {
	output := make(chan generics.KeyValuePair[K, V])
	go func() {
		t.lock.RLock()
		defer func() {
			t.lock.RUnlock()
			close(output)
		}()

		current := t.Root
		for current != nil {
			if current.Children != nil {
				current = current.Children[0]
			} else {
				for i := 0; i < current.Count; i++ {
					output <- generics.KeyValuePair[K, V]{
						Key:   current.Keys[i],
						Value: current.Records[i].Value,
					}
				}
				current = current.NextSibling
			}
		}
	}()

	return output
}
