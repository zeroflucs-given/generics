package ringbuffer

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	collections "github.com/zeroflucs-given/generics/collections"
)

// TestRingBufferInserts tests that we can fill a buffer up to its limit and then
// error
func TestRingBufferInserts(t *testing.T) {
	src := rand.NewSource(133713371337)
	rng := rand.New(src)
	testSize := 13

	buff := New[int](testSize)

	for i := 0; i < testSize; i++ {
		value := rng.Int()
		err := buff.Push(value)
		require.NoError(t, err, "Should not error inserting value %d", i)
	}

	errOverflow := buff.Push(1)
	require.ErrorIs(t, errOverflow, collections.ErrBufferFull, "Should overflow on the last value")
}

// TestRingBufferBasic performs a sequence of operations against the buffer, randomly adding and removing items
// and checking expected conditions.
func TestRingBufferBasic(t *testing.T) {
	src := rand.NewSource(133713371337)
	rng := rand.New(src)
	testSize := 13

	buff := New[int](testSize)
	balance := 0

	require.Equal(t, testSize, buff.Capacity())

	for i := 0; i < 10000; i++ {
		count := buff.Count()
		require.Equal(t, balance, count, "Balance should match count")

		push := rng.Float64() >= 0.5
		if push {
			expectError := balance == testSize
			err := buff.Push(rng.Int())
			gotErr := err != nil
			if gotErr != expectError {
				t.Errorf("Pushing error state: Err=%v, ExpectedErr=%v", err, expectError)
				t.Errorf("Buffer state: %v", buff)
				t.FailNow()
				return
			}

			if err == nil {
				balance++
			}
		} else {
			expectValue := balance > 0
			hasValue, _ := buff.Pop()
			if hasValue != expectValue {
				t.Errorf("Pop value mismatch: HasValue: %v, ExpectValue: %v", hasValue, expectValue)
				t.Errorf("Buffer state: %v", buff)
				t.FailNow()
				return
			}

			if hasValue {
				balance--
			}
		}
	}
}

func TestRingBufferPeek(t *testing.T) {
	buff := New[int](32)

	hasValue, value := buff.Peek()
	require.False(t, hasValue, "Should not have a value to start")
	require.Zero(t, value, "Should have the zero-value for T at start")

	err := buff.Push(4)
	require.NoError(t, err, "Should not error")

	valueRequired, value := buff.Peek()
	require.True(t, valueRequired, "Should have a value to peek")
	require.Equal(t, value, 4, "Both values should be same")
	require.Equal(t, buff.Count(), 1, "Should have a count of one")

}

// BenchmarkRingBuffer tests how fast we can cycle data through the ring-buffer
func BenchmarkRingBuffer(b *testing.B) {
	buff := New[int](32)
	for i := 0; i < buff.Capacity()/2; i++ {
		err := buff.Push(i)
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
	}

	for i := 0; i < b.N; i++ {
		// Put one in
		err := buff.Push(i)
		if err != nil {
			b.Log(err)
			b.FailNow()
		}

		// Pop one
		hasItem, _ := buff.Pop()
		if !hasItem {
			b.Log("Should have had item")
			b.FailNow()
		}
	}

}
