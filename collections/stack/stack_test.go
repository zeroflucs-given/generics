package stack

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	collections "github.com/zeroflucs-given/generics/collections"
)

func TestStackBasics(t *testing.T) {
	stack := NewStack[int](20)

	require.Equal(t, 20, stack.Capacity(), "Should have right capacity")

	for i := 0; i < 20; i++ {
		errPush := stack.Push(i)
		require.NoError(t, errPush, "Should not error")
	}

	errTooMany := stack.Push(1337)
	require.ErrorIs(t, errTooMany, collections.ErrBufferFull, "Should have a buffer overflow error")

	for i := 19; i >= 0; i-- {
		require.Equal(t, i+1, stack.Count(), "Should have the right count")
		found, v := stack.Pop()
		require.True(t, found, "Should have popped value: %v but stack was empty", v)
		require.Equal(t, i, v)
	}

	found, _ := stack.Pop()
	require.False(t, found, "Should have no more values")
}

func TestStackPeek(t *testing.T) {
	buff := NewStack[int](32)

	hasValue, value := buff.Peek()
	require.False(t, hasValue, "Should not have a value to start")
	require.Zero(t, value, "Should have the zero-value for T at start")

	err := buff.Push(4)
	require.NoError(t, err, "Should not error")

	valueRequired, value := buff.Peek()
	require.True(t, valueRequired, "Should have a value to peek")
	require.Equal(t, 4, value, "Both values should be same")
	require.Equal(t, 1, buff.Count(), "Should have a count of one")

}

func BenchmarkStackPushPop(b *testing.B) {
	b.StopTimer()
	dataSize := 10000
	insertP := rand.Perm(dataSize)
	b.StartTimer()
	i := 0

	for i < b.N {
		// Make a stack and push values onto it
		stack := NewStack[int](dataSize)
		for _, value := range insertP {
			err := stack.Push(value)
			if err != nil {
				b.Log(err)
				b.FailNow()
			}

			i++
			if i >= b.N {
				break
			}
		}

		// Clear the stack
		for {
			more, _ := stack.Pop()
			if !more {
				break
			}
		}

		if i >= b.N {
			return
		}
	}
}
