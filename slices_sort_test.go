package generics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSort checks a basic ascending sort
func TestSort(t *testing.T) {
	// Arrange
	in := []int{10, 7, 5, 6, 5, 3, 11, 9}

	// Act
	out := Sort(in)

	// Assert
	require.Equal(t, []int{
		3, 5, 5, 6, 7, 9, 10, 11,
	}, out, "Sorted sequence should be correct")
}

// TestSortDescending checks a basic descending sort
func TestSortDescending(t *testing.T) {
	// Arrange
	in := []int{10, 7, 5, 6, 5, 3, 11, 9}

	// Act
	out := SortDescending(in)

	// Assert
	require.Equal(t, []int{
		11, 10, 9, 7, 6, 5, 5, 3,
	}, out, "Sorted sequence should be correct")
}
