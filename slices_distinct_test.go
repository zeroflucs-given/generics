package generics

import (
	"hash/fnv"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceDistinct(t *testing.T) {
	type ComplexKey struct {
		A string
		B string
	}

	keys := []ComplexKey{
		{A: "b", B: "a"},
		{A: "c", B: "c"},
		{A: "b", B: "a"},
		{A: "c", B: "b"},
		{A: "c", B: "b"},
		{A: "c", B: "b"},
		{A: "c", B: "b"},
		{A: "a", B: "a"},
		{A: "a", B: "c"},
		{A: "a", B: "c"},
		{A: "a", B: "c"},
		{A: "a", B: "b"},
	}

	t.Run("Distinct", func(t *testing.T) {
		simpleKeys := Map(keys, func(index int, input ComplexKey) string { return input.A })

		xx := Distinct(simpleKeys)

		expected := []string{"a", "b", "c"}

		assert.Equal(t, expected, xx)
	})

	t.Run("DistinctFunc", func(t *testing.T) {
		xx := DistinctFunc(slices.Clone(keys), func(a, b ComplexKey) int {
			if val := strings.Compare(a.A, b.A); val != 0 {
				return val
			}

			return strings.Compare(a.B, b.B)
		})

		expected := []ComplexKey{
			{A: "a", B: "a"},
			{A: "a", B: "b"},
			{A: "a", B: "c"},
			{A: "b", B: "a"},
			{A: "c", B: "b"},
			{A: "c", B: "c"},
		}

		assert.Equal(t, expected, xx)
	})

	t.Run("DistinctStable", func(t *testing.T) {
		xx := DistinctStable(slices.Clone(keys))

		expected := []ComplexKey{
			{A: "b", B: "a"},
			{A: "c", B: "c"},
			{A: "c", B: "b"},
			{A: "a", B: "a"},
			{A: "a", B: "c"},
			{A: "a", B: "b"},
		}

		assert.Equal(t, expected, xx)
	})

	t.Run("DistinctStableFunc", func(t *testing.T) {
		xx := DistinctStableFunc(slices.Clone(keys), func(val ComplexKey) uint64 {
			h := fnv.New64a()
			_, _ = h.Write([]byte(val.A + val.B))
			return h.Sum64()
		})

		expected := []ComplexKey{
			{A: "b", B: "a"},
			{A: "c", B: "c"},
			{A: "c", B: "b"},
			{A: "a", B: "a"},
			{A: "a", B: "c"},
			{A: "a", B: "b"},
		}

		assert.Equal(t, expected, xx)
	})
}
