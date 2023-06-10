package main

import (
	"fmt"
)

func combinations(objects []interface{}, n int) [][]interface{} {
	if n == 0 || len(objects) == 0 {
		return [][]interface{}{}
	}

	combs := [][]interface{}{}
	for i, obj := range objects {
		if n == 1 {
			combs = append(combs, []interface{}{obj})
		} else {
			rest := objects[i+1:]
			for _, comb := range combinations(rest, n-1) {
				combs = append(combs, append([]interface{}{obj}, comb...))
			}
		}
	}

	return combs
}

func main() {
	objects := []interface{}{"a", "b", "c", "d", "e"}
	n := 3
	result := combinations(objects, n)
	fmt.Printf("Combinations of size %d: %v\n", n, result)
}
