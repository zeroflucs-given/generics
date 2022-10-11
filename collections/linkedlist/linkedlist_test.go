package linkedlist

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeroflucs-given/generics/collections"
)

func TestLinkedListPeek(t *testing.T) {
	buff := New[int]()

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

func TestLinkedListBasic(t *testing.T) {
	testSize := 10000

	buff := New[int]()
	require.Equal(t, collections.CapacityInfinite, buff.Capacity())

	for i := 0; i < testSize; i++ {
		err := buff.Push(i)
		require.NoError(t, err, "Should not error inserting value %d", i)
	}

	for i := 0; i < testSize; i++ {
		hasValue, value := buff.Pop()
		require.True(t, hasValue, "Should have a value to pop")
		require.Equal(t, i, value, "Should have the same value")
	}

}

func TestLinkedListContains(t *testing.T) {
	buff := New[int]()
	require.Equal(t, collections.CapacityInfinite, buff.Capacity())

	err := buff.Insert(45)
	require.NoErrorf(t, err, "Should not error inserting")

	hasValue := buff.Contains(45)
	require.True(t, hasValue, "Should be contain the value")

	buff.Remove(45)
	require.False(t, buff.Contains(45), "should no longer contain anything")
}

func TestLinkedListRemoveFirst(t *testing.T) {
	buff := New[int]()
	require.Equal(t, collections.CapacityInfinite, buff.Capacity())

	for i := 0; i < 5; i++ {
		err := buff.Insert(i)
		require.NoErrorf(t, err, "Should not error inserting")
	}

	buff.Remove(0)
	require.False(t, buff.Contains(0), "Should not have 0 in the list")
	require.Equal(t, buff.Count(), 4, "Should have 4 items in the list")
}

func TestLinkedListRemoveLast(t *testing.T) {
	buff := New[int]()
	require.Equal(t, collections.CapacityInfinite, buff.Capacity())

	for i := 0; i < 5; i++ {
		err := buff.Insert(i)
		require.NoErrorf(t, err, "Should not error inserting")
	}

	buff.Remove(1)
	require.False(t, buff.Contains(1), "Should not have 4 in the list")
	require.Equal(t, buff.Count(), 4, "Should have 4 items in the list")
}

func TestLinkedListRemoveFirstIndex(t *testing.T) {
	buff := New[int]()
	require.Equal(t, collections.CapacityInfinite, buff.Capacity())

	for i := 0; i < 5; i++ {
		err := buff.Insert(i)
		require.NoErrorf(t, err, "Should not error inserting")
	}

	buff.RemoveAt(4)

	require.False(t, buff.Contains(0), "Should not have 4 in the list")
	require.Equal(t, 4, buff.Count(), "Should have 4 items in the list")
}

func TestLinkedListRemoveLastIndex(t *testing.T) {
	buff := New[int]()
	require.Equal(t, collections.CapacityInfinite, buff.Capacity())

	for i := 0; i < 5; i++ {
		err := buff.Insert(i)
		require.NoErrorf(t, err, "Should not error inserting")
	}

	buff.RemoveAt(0)

	require.False(t, buff.Contains(4), "Should not have 4 in the list")
	require.Equal(t, 4, buff.Count(), "Should have 4 items in the list")
}

func TestLinkedListRemoveIndexOutOfBounds(t *testing.T) {
	buff := New[int]()
	require.Equal(t, collections.CapacityInfinite, buff.Capacity())

	for i := 0; i < 5; i++ {
		err := buff.Insert(i)
		require.NoErrorf(t, err, "Should not error inserting")
	}

	require.Panics(t, func() {
		buff.RemoveAt(5)
	})
}

func TestLinkedListValue(t *testing.T) {
	buff := New[int]()
	require.Equal(t, collections.CapacityInfinite, buff.Capacity())

	for i := 0; i < 5; i++ {
		err := buff.Insert(i)
		require.NoErrorf(t, err, "Should not error inserting")
	}

	hasValue, value := buff.Value(0)
	require.Equal(t, 4, value, "Both values should be same")
	require.True(t, hasValue, "Should have the value 4")
}

func TestLinkedListValueAtHead(t *testing.T) {
	buff := New[int]()
	require.Equal(t, collections.CapacityInfinite, buff.Capacity())

	for i := 0; i < 5; i++ {
		err := buff.Insert(i)
		require.NoErrorf(t, err, "Should not error inserting")
	}

	hasValue, value := buff.Value(4)
	require.Equal(t, 0, value, "Both values should be same")
	require.True(t, hasValue, "Should have the value 4")
}
