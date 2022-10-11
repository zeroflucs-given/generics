package bplustree

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkTreeNodeGetInsertIndex(b *testing.B) {
	for _, v := range treeOrders {
		order := v
		b.Run(fmt.Sprintf("TreeOrder=%d Linear", order), func(b *testing.B) {
			benchmarkTreeNodeGetInsertIndexInternal(b, order)
		})
	}
}

func benchmarkTreeNodeGetInsertIndexInternal(b *testing.B, order int) {
	src := rand.NewSource(int64(order * 13371337)) // Constant value
	rnd := rand.New(src)

	node := treeNode[int32, int]{
		Count: order,
		Keys:  make([]int32, order),
	}

	// Each cycle, use some new values
	b.StopTimer()
	for initIndex := 0; initIndex < order; initIndex++ {
		node.Keys[initIndex] = rnd.Int31()
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		valueToInsert := rnd.Int31()
		node.getInsertIndex(valueToInsert)
	}
}
