package generics_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeroflucs-given/generics"
)

func TestConcatenate(t *testing.T) {
	// Arrange
	first := []int{1, 2, 3}
	second := []int{4, 5, 6}
	third := []int{7, 8, 9}

	// Act
	result := generics.Concatenate(first, second, third)

	// Assert
	require.NotNil(t, result, "Should have an object")
	require.Len(t, result, 9, "Should have right number of items")
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, result, "Should have right result")
}

func TestContainsTrue(t *testing.T) {
	// Arrange
	input := []int{10, 14, 28}
	target := 14

	// Act
	result := generics.Contains(input, target)

	// Assert
	require.True(t, result, "Should have the correct result")
}

func TestContainsFalse(t *testing.T) {
	// Arrange
	input := []int{10, 14, 28}
	target := 43

	// Act
	result := generics.Contains(input, target)

	// Assert
	require.False(t, result, "Should have the correct result")
}

func TestCut(t *testing.T) {
	// Arrange
	input := []int{9, 8, 7, 6}

	// Act
	head, rest := generics.Cut(input)

	// Assert
	assert.Equal(t, 9, head, "Should have the right head value")
	assert.Equal(t, []int{8, 7, 6}, rest, "Should have correct list residual")
}

func TestCutLast(t *testing.T) {
	// Arrange
	input := []int{9}

	// Act
	head, rest := generics.Cut(input)

	// Assert
	assert.Equal(t, 9, head, "Should have the right head value")
	assert.Nil(t, rest, "Should have no list residual")
}

func TestCutEmpty(t *testing.T) {
	// Arrange
	var input []int

	// Act
	head, rest := generics.Cut(input)

	// Assert
	assert.Equal(t, 0, head, "Should have correct head (type default)")
	assert.Nil(t, rest, "Empty slice")
}

func TestDefaultIfEmpty(t *testing.T) {
	// Arrange
	input := []int{}
	def := 42

	// Act
	result := generics.DefaultIfEmpty(input, def)

	// Assert
	require.Len(t, result, 1, "Should have right length")
	require.Equal(t, result[0], def, "Should have the right value")
}

func TestDefaultIfEmptyNotUsed(t *testing.T) {
	// Arrange
	input := []int{1, 3, 7}
	def := 42

	// Act
	result := generics.DefaultIfEmpty(input, def)

	// Assert
	require.Len(t, result, 3, "Should have right length")
	require.Equal(t, input, result, "Should have the right value")
}

func TestReverse(t *testing.T) {
	// Arrange
	input := []int{5, 4, 3}

	// Act
	result := generics.Reverse(input)

	// Assert
	require.NotNil(t, result, "Should have a result")
	require.Equal(t, []int{3, 4, 5}, result, "Should have right result")
}

func TestReverseEmpty(t *testing.T) {
	// Arrange
	input := []int{}

	// Act
	result := generics.Reverse(input)

	// Assert
	require.NotNil(t, result, "Should have a result")
	require.Len(t, result, 0, "Should have the right result")
}

func TestSkip(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4, 5, 6}

	// Act
	result := generics.Skip(input, 3)

	// Assert
	require.NotNil(t, result, "Should have a result")
	require.Equal(t, []int{4, 5, 6}, result, "Should have right result")
}

func TestSkipTooMany(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4, 5, 6}

	// Act
	result := generics.Skip(input, 99)

	// Assert
	require.Nil(t, result, "Should not have a result")
}

func TestTake(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4, 5, 6}

	// Act
	result := generics.Take(input, 3)

	// Assert
	require.NotNil(t, result, "Should have a result")
	require.Equal(t, []int{1, 2, 3}, result, "Should have right result")
}

func TestTakeEmpty(t *testing.T) {
	// Arrange
	input := []int{}

	// Act
	result := generics.Take(input, 99)

	// Assert
	require.Nil(t, result, "Should not have a result")
}

func TestTakeTooMany(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4, 5, 6}

	// Act
	result := generics.Take(input, 99)

	// Assert
	require.NotNil(t, result, "Should not have a result")
	require.Equal(t, input, result, "Should have the right result")
}

func TestTakeWhile(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	filter := func(index int, v int) bool {
		return v < 5
	}

	// Act
	taken := generics.TakeWhile(input, filter)

	// Assert
	require.Equal(t, []int{1, 2, 3, 4}, taken, "Should have the right output")
}

func TestTakeWhileContext(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	testCtx := context.Background()
	filter := func(ctx context.Context, index int, v int) (bool, error) {
		return v < 5, nil
	}

	// Act
	taken, err := generics.TakeWhileWithContext(testCtx, input, filter)

	// Assert
	require.NoError(t, err, "Should not error")
	require.Equal(t, []int{1, 2, 3, 4}, taken, "Should have the right output")
}

func TestTakeWhileContextError(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	testCtx := context.Background()
	filter := func(ctx context.Context, index int, v int) (bool, error) {
		if v > 4 {
			return false, errors.New("boom")
		}
		return true, nil
	}

	// Act
	taken, err := generics.TakeWhileWithContext(testCtx, input, filter)

	// Assert
	require.Error(t, err, "Should propegate error")
	require.Nil(t, taken, "Should have no output")
}

func TestTakeUntil(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	filter := func(index int, v int) bool {
		return v > 4
	}

	// Act
	taken := generics.TakeUntil(input, filter)

	// Assert
	require.Equal(t, []int{1, 2, 3, 4}, taken, "Should have the right output")
}

func TestTakeUntilContext(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	testCtx := context.Background()
	filter := func(ctx context.Context, index int, v int) (bool, error) {
		return v > 4, nil
	}

	// Act
	taken, err := generics.TakeUntilWithContext(testCtx, input, filter)

	// Assert
	require.NoError(t, err, "Should not error")
	require.Equal(t, []int{1, 2, 3, 4}, taken, "Should have the right output")
}

func TestTakeUntilContextError(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	testCtx := context.Background()
	filter := func(ctx context.Context, index int, v int) (bool, error) {
		if v > 4 {
			return false, errors.New("boom")
		}
		return false, nil
	}

	// Act
	taken, err := generics.TakeUntilWithContext(testCtx, input, filter)

	// Assert
	require.Error(t, err, "Should propegate error")
	require.Nil(t, taken, "Should have no output")
}
