package weightedrandom

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestWeightedRandomBasics cecks that we do some weighting of values accoring to scores
func TestWeightedRandomBasics(t *testing.T) {
	initialWeight := 1000.0
	items := 20
	samples := 1000 * 1000 * 10 // 10M

	src := rand.NewSource(133713371337)
	rng := rand.New(src)

	// Push the values in
	wr := NewWeightedRandom[int](items, UnsignedFloatMapper)

	hasPick, v := wr.Pick(rng)
	require.False(t, hasPick, "Should not be able to pick when empty")
	require.Zero(t, v, "Should have no value from the picker")

	currentWeight := initialWeight
	for i := 0; i < items; i++ {
		itemWeight := currentWeight
		currentWeight /= 2.0
		wr.Push(itemWeight, i)
	}

	frequencies := make([]int, items)
	for i := 0; i < samples; i++ {
		found, picked := wr.Pick(rng)
		require.True(t, found, "Should find a value")
		frequencies[picked]++
	}

	for i := 0; i < items; i++ {
		freq := frequencies[i]
		max := (float64(samples) / math.Pow(2, float64(i+1)))
		maxWithTolerance := max * 1.1

		t.Logf("Item %v hits %v (%.2f%%) [Expect ~%d]", i, freq, 100.0*float64(freq)/float64(samples), int(max))

		if freq > int(maxWithTolerance) && freq > 100 {
			t.Logf("Too many hits at this frequency. Should have a maximum of %v", max)
		}
	}
}

func TestWeightedRandomDefectiveMapper(t *testing.T) {
	require.Panics(t, func() {
		wr := NewWeightedRandom[int](1, func(v int) float64 {
			return -1
		})
		wr.Push(42, 1337)
	}, "This code should panic")
}

// BenchmarkWeightedRandomPick does weighted random selection using a series of values
// to see how fast we can hit the code.
func BenchmarkWeightedRandomPick(b *testing.B) {
	for i := 1; i < 8; i++ {
		items := int(math.Pow(10, float64(i)))

		b.Run(fmt.Sprintf("Count=%d", items), func(b *testing.B) {
			b.StopTimer()

			src := rand.NewSource(133713371337)
			rng := rand.New(src)
			wr := NewWeightedRandom[int](items, UnsignedFloatMapper)
			for i := 0; i < items; i++ {
				wr.Push(float64(rng.Int31n(10000)), i)
			}

			b.StartTimer()

			for i := 0; i < b.N; i++ {
				wr.Pick(rng)
			}
		})
	}
}
