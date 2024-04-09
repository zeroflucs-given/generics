package generics_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeroflucs-given/generics"
	"testing"
)

func TestChunk(t *testing.T) {
	items := []int{
		1, 1, 2, 3, 5, 8, 13, 21,
	}

	type testDef struct {
		Buckets  int
		Expected [][]int
	}

	tests := []testDef{
		{
			Buckets:  0,
			Expected: nil,
		},
		{
			Buckets:  1,
			Expected: [][]int{{1, 1, 2, 3, 5, 8, 13, 21}},
		},
		{
			Buckets:  2,
			Expected: [][]int{{1, 1, 2, 3}, {5, 8, 13, 21}},
		},
		{
			Buckets:  3,
			Expected: [][]int{{1, 1, 2}, {3, 5, 8}, {13, 21}},
		},
		{
			Buckets:  4,
			Expected: [][]int{{1, 1}, {2, 3}, {5, 8}, {13, 21}},
		},
		{
			Buckets:  5,
			Expected: [][]int{{1, 1}, {2, 3}, {5, 8}, {13}, {21}},
		},
		{
			Buckets:  6,
			Expected: [][]int{{1, 1}, {2, 3}, {5}, {8}, {13}, {21}},
		},
		{
			Buckets:  7,
			Expected: [][]int{{1, 1}, {2}, {3}, {5}, {8}, {13}, {21}},
		},
		{
			Buckets:  8,
			Expected: [][]int{{1}, {1}, {2}, {3}, {5}, {8}, {13}, {21}},
		},
		{
			Buckets:  9,
			Expected: [][]int{{1}, {1}, {2}, {3}, {5}, {8}, {13}, {21}, {}},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.Buckets), func(t *testing.T) {
			r := generics.Chunk(items, test.Buckets)
			require.Equal(t, test.Expected, r)
		})
	}
}

var benchArray []int

func init() {
	benchArray = make([]int, 1_000_000_000)
	for i := range benchArray {
		benchArray[i] = i
	}
}

func BenchmarkChunk(b *testing.B) {
	_ = generics.Chunk(benchArray, b.N)

	b.StopTimer()

	b.Logf("%v buckets: %v", b.N, b.Elapsed())
}

func TestChunkMap(t *testing.T) {
	items := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8}

	type testDef struct {
		ExpectedCounts []int
	}

	tests := []testDef{
		{ExpectedCounts: []int{}},
		{ExpectedCounts: []int{8}},
		{ExpectedCounts: []int{4, 4}},
		{ExpectedCounts: []int{3, 3, 2}},
		{ExpectedCounts: []int{2, 2, 2, 2}},
		{ExpectedCounts: []int{2, 2, 2, 1, 1}},
		{ExpectedCounts: []int{2, 2, 1, 1, 1, 1}},
		{ExpectedCounts: []int{2, 1, 1, 1, 1, 1, 1}},
		{ExpectedCounts: []int{1, 1, 1, 1, 1, 1, 1, 1}},
		{ExpectedCounts: []int{1, 1, 1, 1, 1, 1, 1, 1, 0}},
	}

	// Map iteration order is non-deterministic, so combine them all and make sure
	// we have the input.
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", len(test.ExpectedCounts)), func(t *testing.T) {
			r := generics.ChunkMap(items, len(test.ExpectedCounts))
			require.Len(t, r, len(test.ExpectedCounts))

			comb := make(map[int]int, len(items))
			for i, m := range r {
				assert.Len(t, m, test.ExpectedCounts[i])
				for k, v := range m {
					comb[k] = v
				}
			}

			if len(r) > 0 {
				require.Equal(t, items, comb)
			}
		})
	}
}

var chunkBenchMap map[int]int

func init() {
	chunkBenchMap = make(map[int]int, 1_000_000_000)
	for i := range chunkBenchMap {
		chunkBenchMap[i] = i
	}
}

func BenchmarkChunkMap(b *testing.B) {
	_ = generics.ChunkMap(chunkBenchMap, b.N)

	b.StopTimer()

	b.Logf("%v buckets: %v", b.N, b.Elapsed())
}
