package query_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeroflucs-given/generics/query"
)

func ExampleSlice() {
	inputs := []int{
		1, 2, 3, 4, 5,
		6, 7, 8, 9, 10,
		11, 12, 13, 14, 15,
		16, 17, 18, 19, 20,
	}
	query := query.Slice(inputs)

	result := query.
		Skip(5).
		TakeUntil(func(index int, v int) bool {
			return v >= 15
		}).
		Filter(func(index int, v int) bool {
			return v%2 == 0 // Evens only
		}).
		Reverse().
		ToSlice()

	fmt.Println(result) // [14, 12, 10, 8, 6]
}

// TestQueryableSlice just checks we can chain up long sequences of mutations
func TestQueryableSlice(t *testing.T) {
	// Arrange
	inputs := [][]int{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10},
		{11, 12, 13, 14, 15},
		{16, 17, 18, 19, 20},
	}
	// Act
	query := query.Slice(inputs[0]).
		Concatenate(inputs[1:]...).
		TakeWhile(func(index int, v int) bool {
			return index < 18
		}).
		Skip(5).
		TakeUntil(func(index int, v int) bool {
			return v >= 15
		}).
		Filter(func(index int, v int) bool {
			return v%2 == 0 // Evens only
		}).
		Mutate(func(index int, v int) int {
			return v + 2
		}).
		Take(3).
		Reverse()

	all := query.All(func(index int, v int) bool {
		return v < 13
	})
	any := query.Any(func(index int, v int) bool {
		return v%2 == 1
	})
	count := query.Count()
	result := query.ToSlice()
	first := query.First()
	last := query.Last()

	require.True(t, all, "All items should be < 12")
	require.False(t, any, "Should have no odd values")
	require.Equal(t, 3, count, "Should have 3 items")
	require.Equal(t, 12, first, "First item should be 10")
	require.Equal(t, 8, last, "Last item should be 6")
	require.Equal(t, []int{12, 10, 8}, result, "Should have right items")
}
