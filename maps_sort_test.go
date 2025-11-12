package generics_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeroflucs-given/generics"
)

// TestSortedByKey checks we get keys in order
func TestSortedByKey(t *testing.T) {
	input := map[int]string{
		2: "Hey, ",
		4: "World",
		3: "Hello",
	}

	sorted := generics.SortedByKey(input)

	require.Len(t, sorted, 3, "Should have 3 items")
	require.Equal(t, "Hey, ", sorted[0].Value)
	require.Equal(t, "Hello", sorted[1].Value)
	require.Equal(t, "World", sorted[2].Value)
}
