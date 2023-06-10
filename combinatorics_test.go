package generics_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeroflucs-given/generics"
)

func TestCombinations(t *testing.T) {
	// Arrange
	input := []string{"a", "b", "c", "d"}

	// Act
	output := generics.Combinations(input, 2)

	// Assert
	assert.Len(t, output, 6, "Should have correct number of combinations")
}

func TestCombinationsSize(t *testing.T) {
	// Arrange
	input := []string{"a", "b", "c"}

	// Act
	output := generics.Combinations(input, 4)

	// Assert
	require.Nil(t, output)
}

func TestCombinationsFiltered(t *testing.T) {
	// Arrange
	input := []string{"a", "b", "c", "d"}

	// Act
	output := generics.CombinationsFiltered(input, 2, func(i int, v string) bool {
		return i%2 == 0 // Only include even indexes
	})

	// Assert
	assert.Len(t, output, 1, "Should have correct number of combinations")
	assert.Equal(t, generics.IndexedItem[string]{Index: 0, Item: "a"}, output[0][0])
	assert.Equal(t, generics.IndexedItem[string]{Index: 2, Item: "c"}, output[0][1])
}
