//go:build perfanalysis
// +build perfanalysis

package bplustree

import (
	"fmt"
	"math/rand"
	"testing"
)

var dataSizes = []int{10000, 100000, 250000, 500000}

// BenchmarkRandomInserts checks how fast we can perform random inserts
// into the tree.
func BenchmarkRandomInsertPermutations(b *testing.B) {
	skipOthers := false

	for _, v := range treeOrders {
		order := v

		for _, v := range dataSizes {
			dataSize := v

			b.Run(fmt.Sprintf("Insert_Order=%d_Count=%d_Random", order, dataSize), func(b *testing.B) {
				b.ResetTimer()

				if skipOthers {
					b.SkipNow()
				}

				benchmarkInsertsRandomInternal(b, order, dataSize)
			})
		}
	}
}

// BenchmarkInsertsRandomFixedMedium performs a fixed test with order 27, and 0.5m
// data points.
func BenchmarkInsertsRandomFixedMedium(b *testing.B) {
	benchmarkInsertsRandomInternal(b, 27, 500000)
}

// BenchmarkInsertsRandomFixedLarge performs a fixed test with order 27, and 2.5m
// data points.
func BenchmarkInsertsRandomFixedLarge(b *testing.B) {
	benchmarkInsertsRandomInternal(b, 27, 2500000)
}

// BenchmarkInsertsSequentialFixedMedium performs a fixed test with order 27, and 0.5m
// data points using linear values
func BenchmarkInsertsSequentialFixedMedium(b *testing.B) {
	benchmarkInsertsSequentialInternal(b, 27, 500000)
}

// BenchmarkInsertsSequentialFixedLarge performs a fixed test with order 27, and 2.5m
// data points using linear values
func BenchmarkInsertsSequentialFixedLarge(b *testing.B) {
	benchmarkInsertsSequentialInternal(b, 27, 2500000)
}

func benchmarkInsertsRandomInternal(b *testing.B, order int, dataSize int) {
	for iter := 0; iter < b.N; iter++ {
		tree, err := New[int, int](order, DefaultTestPreAlloc)
		if err != nil {
			b.Logf("Error: %v", err)
			b.FailNow()
		}

		rs := rand.NewSource(int64(order * dataSize))
		rnd := rand.New(rs)

		for i := 0; i < dataSize; i++ {
			k := int(rnd.Int31n(1000000))
			v := i

			tree.Insert(k, v)
		}
	}
}

func benchmarkInsertsSequentialInternal(b *testing.B, order int, dataSize int) {
	for iter := 0; iter < b.N; iter++ {
		tree, err := New[int, int](order, DefaultTestPreAlloc)
		if err != nil {
			b.Logf("Error: %v", err)
			b.FailNow()
		}

		for i := 0; i < dataSize; i++ {
			k := i
			v := i

			tree.Insert(k, v)
		}
	}
}
