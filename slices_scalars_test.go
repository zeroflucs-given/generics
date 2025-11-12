package generics_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeroflucs-given/generics"
)

// TestAll checks a simple case where we return a true value
func TestAll(t *testing.T) {
	// Arrange
	filter := func(index int, i int) bool {
		return i == index+1
	}
	input := []int{1, 2, 3, 4}

	// Act
	result := generics.All(input, filter)

	// Assert
	require.True(t, result, "Should have the correct result")
}

// TestAllFail checks a simple negative case where one value is wrong.
// This also validates we are operating in a lazy fashion.
func TestAllFail(t *testing.T) {
	// Arrange
	calls := 0
	filter := func(index int, i int) bool {
		calls++
		return i == index+1
	}
	input := []int{1, 2, 4, 3}

	// Act
	result := generics.All(input, filter)

	// Assert
	require.False(t, result, "Should have the correct result")
	require.Equal(t, 3, calls, "Should be lazy")
}

// TestAllContext checks a positive outcome with a context
func TestAllContext(t *testing.T) {
	// Arrange
	testCtx := context.TODO()
	filter := func(ctx context.Context, index int, i int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		return i == index+1, nil
	}
	input := []int{1, 2, 3, 4}

	// Act
	result, err := generics.AllWithContext(testCtx, input, filter)

	// Assert
	require.NoError(t, err, "Should not error")
	require.True(t, result, "Should have the correct result")
}

// TestAllContextFail checks a negative case with a context
func TestAllContextFail(t *testing.T) {
	// Arrange
	testCtx := context.TODO()
	calls := 0
	filter := func(ctx context.Context, index int, i int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		calls++
		return i == index+1, nil
	}
	input := []int{1, 2, 4, 3}

	// Act
	result, err := generics.AllWithContext(testCtx, input, filter)

	// Assert
	require.NoError(t, err, "Should not error")
	require.False(t, result, "Should have the correct result")
	require.Equal(t, 3, calls, "Should be lazy")
}

// TestAllContextError checks that All propegates context errors from filters
func TestAllContextError(t *testing.T) {
	// Arrange
	calls := 0
	testCtx := context.TODO()
	filter := func(ctx context.Context, index int, i int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		calls++
		if index == 2 {
			return false, errors.New("failure")
		}
		return i == index+1, nil
	}
	input := []int{1, 2, 3, 4}

	// Act
	result, err := generics.AllWithContext(testCtx, input, filter)

	// Assert
	require.Error(t, err, "Should error")
	require.False(t, result, "Should have the correct result")
	require.Equal(t, 3, calls, "Should be lazy")
}

// TestAny checks a positive case for "any"
func TestAny(t *testing.T) {
	// Arrange
	calls := 0
	filter := func(index int, i int) bool {
		calls++
		return i == 2
	}
	input := []int{1, 2, 3, 4}

	// Act
	result := generics.Any(input, filter)

	// Assert
	require.True(t, result, "Should have the correct result")
	require.Equal(t, 2, calls, "Should be lazy")
}

// TestAnyFail checks a negative case for "any"
func TestAnyFail(t *testing.T) {
	// Arrange
	calls := 0
	filter := func(index int, i int) bool {
		calls++
		return i == 9
	}
	input := []int{1, 2, 4, 3}

	// Act
	result := generics.Any(input, filter)

	// Assert
	require.False(t, result, "Should have the correct result")
	require.Equal(t, len(input), calls, "Should check all items")
}

// TestAnyContext checks a positive case matching any, verifying we are lazy
func TestAnyContext(t *testing.T) {
	// Arrange
	calls := 0
	testCtx := context.TODO()
	filter := func(ctx context.Context, index int, i int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		calls++
		return i == 3, nil
	}
	input := []int{1, 2, 3, 4}

	// Act
	result, err := generics.AnyWithContext(testCtx, input, filter)

	// Assert
	require.NoError(t, err, "Should not error")
	require.True(t, result, "Should have the correct result")
	require.Equal(t, 3, calls, "Should be lazy")
}

// TestAnyContextFail checks a negative case for any, verifying we are greedy
func TestAnyContextFail(t *testing.T) {
	// Arrange
	testCtx := context.TODO()
	calls := 0
	filter := func(ctx context.Context, index int, i int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		calls++
		return i == 9, nil
	}
	input := []int{1, 2, 4, 3}

	// Act
	result, err := generics.AnyWithContext(testCtx, input, filter)

	// Assert
	require.NoError(t, err, "Should not error")
	require.False(t, result, "Should have the correct result")
	require.Equal(t, len(input), calls, "Should check all items")
}

// TestAnyContextError checks Any propegates context errors
func TestAnyContextError(t *testing.T) {
	// Arrange
	calls := 0
	testCtx := context.TODO()
	filter := func(ctx context.Context, index int, i int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		calls++
		if index == 2 {
			return false, errors.New("failure")
		}
		return i == 9, nil // We won't get here
	}
	input := []int{1, 2, 3, 4}

	// Act
	result, err := generics.AnyWithContext(testCtx, input, filter)

	// Assert
	require.Error(t, err, "Should error")
	require.False(t, result, "Should have the correct result")
	require.Equal(t, 3, calls, "Should be lazy")
}

// TestCount checks a positive case for "any"
func TestCount(t *testing.T) {
	// Arrange
	calls := 0
	filter := func(index int, i int) bool {
		calls++
		return i%2 == 0
	}
	input := []int{1, 2, 3, 4}

	// Act
	result := generics.Count(input, filter)

	// Assert
	require.Equal(t, 2, result, "Should have correct result")
	require.Equal(t, len(input), calls, "Should be eager")
}

// TestCountNoFilter checks we handle an optimized case where we dont have any filters
func TestCountNoFilter(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4}

	// Act
	result := generics.Count(input)

	// Assert
	require.Equal(t, 4, result, "Should have correct result")
}

// TestCountContext perfoms counting with a filter
func TestCountContext(t *testing.T) {
	// Arrange
	calls := 0
	testCtx := context.TODO()
	filter := func(ctx context.Context, index int, i int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		calls++
		return i%2 == 0, nil
	}
	input := []int{1, 2, 3, 4}

	// Act
	result, err := generics.CountWithContext(testCtx, input, filter)

	// Assert
	require.NoError(t, err, "Should not error")
	require.Equal(t, 2, result, "Should have correct result")
	require.Equal(t, len(input), calls, "Should be eager")
}

// TestCountContextNoFilters checks we count with no filters
func TestCountContextNoFilters(t *testing.T) {
	// Arrange
	testCtx := context.TODO()
	input := []int{1, 2, 3, 4}

	// Act
	result, err := generics.CountWithContext(testCtx, input)

	// Assert
	require.NoError(t, err, "Should not error")
	require.Equal(t, 4, result, "Should have correct result")
}

// TestCountContextError checks we handle context failures
func TestCountContextError(t *testing.T) {
	// Arrange
	calls := 0
	testCtx := context.TODO()
	filter := func(ctx context.Context, index int, i int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		calls++
		if index == 2 {
			return false, errors.New("failure")
		}
		return i%2 == 0, nil
	}
	input := []int{1, 2, 3, 4}

	// Act
	result, err := generics.CountWithContext(testCtx, input, filter)

	// Assert
	require.Error(t, err, "Should error")
	require.Zero(t, result, "Should return default value")
	require.Equal(t, 3, calls, "Should be lazy")
}
