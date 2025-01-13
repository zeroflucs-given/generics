package generics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCoalesce checks the coalesce does what we expect
func TestCoalesce(t *testing.T) {
	input := []string{"", "", "Foo", "Bar"}

	result := Coalesce(input...)

	require.Equal(t, "Foo", result, "Should should give the first non-empty result")
}

// TestCoalesceEmpty checks the coalesce does what we expect when the input is empty
func TestCoalesceEmpty(t *testing.T) {
	input := []int{}

	result := Coalesce(input...)

	require.Equal(t, 0, result, "Should should give the first non-empty result")
}
