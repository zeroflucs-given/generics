package generics_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeroflucs-given/generics"
)

func TestRepeat(t *testing.T) {
	value := 42
	target := 5

	result := generics.Repeat(value, target)

	require.Len(t, result, target)
	require.ElementsMatch(t, result, []int{value, value, value, value, value})
}
