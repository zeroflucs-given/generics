package generics_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeroflucs-given/generics"
)

func TestFirst(t *testing.T) {
	// Arrange
	input := []int{10, 14, 28}
	target := 10

	// Act
	result := generics.First(input)

	// Assert
	require.Equal(t, result, target, "Should have the correct value")
}

func TestFirstEmpty(t *testing.T) {
	// Arrange
	input := []int{}

	// Act
	result := generics.First(input)

	// Assert
	require.Equal(t, 0, result, "Should get a default back")
}

func TestFirstWithContext(t *testing.T) {
	// Arrange
	testCtx := context.Background()
	input := []int{10, 14, 28}
	filter := func(ctx context.Context, index int, v int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		if v >= 12 {
			return true, nil
		}
		return false, nil
	}
	target := 14

	// Act
	result, err := generics.FirstWithContext(testCtx, input, filter)

	// Assert
	require.NoError(t, err, "Should not error")
	require.Equal(t, result, target, "Should have the correct value")
}

func TestFirstWithContextNotFound(t *testing.T) {
	// Arrange
	testCtx := context.Background()
	input := []int{10, 14, 28}
	filter := func(ctx context.Context, index int, v int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		return false, nil
	}
	target := 0

	// Act
	result, err := generics.FirstWithContext(testCtx, input, filter)

	// Assert
	require.NoError(t, err, "Should not error")
	require.Equal(t, result, target, "Should have the correct value")
}

func TestFirstWithContextError(t *testing.T) {
	// Arrange
	testCtx := context.Background()
	input := []int{10, 14, 28}
	filter := func(ctx context.Context, index int, v int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		if v >= 12 {
			return false, errors.New("boom")
		}
		return false, nil
	}

	// Act
	result, err := generics.FirstWithContext(testCtx, input, filter)

	// Assert
	require.Error(t, err, "Should not error")
	require.Equal(t, result, 0, "Should have the default type value")
}

// TestFilter applies a filter to a set
func TestFilter(t *testing.T) {
	// Arrange
	filter := func(index int, v int) bool {
		return v%2 == 0
	}
	input := []int{1, 2, 3, 4, 5, 6}

	// Act
	result := generics.Filter(input, filter)

	// Assert
	require.Len(t, result, 3, "Should have 3 items")
	require.Equal(t, []int{2, 4, 6}, result, "Should have the right outputs")
}

// TestFilterContext applies a filter to a set with a context
func TestFilterContext(t *testing.T) {
	// Arrange
	testCtx := context.TODO()
	filter := func(ctx context.Context, index int, v int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should get the same context back")
		return v%2 == 0, nil
	}
	input := []int{1, 2, 3, 4, 5, 6}

	// Act
	result, err := generics.FilterWithContext(testCtx, input, filter)

	// Assert
	require.NoError(t, err, "Should not error")
	require.Len(t, result, 3, "Should have 3 items")
	require.Equal(t, []int{2, 4, 6}, result, "Should have the right outputs")
}

// TestFilterContextError make sure we abort when errors occur
func TestFilterContextError(t *testing.T) {
	// Arrange
	testCtx := context.TODO()
	filter := func(ctx context.Context, index int, v int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should get the same context back")
		if v == 5 {
			return false, errors.New("boom")
		}
		return v%2 == 0, nil
	}
	input := []int{1, 2, 3, 4, 5, 6}

	// Act
	result, err := generics.FilterWithContext(testCtx, input, filter)

	// Assert
	require.Error(t, err, "Should error")
	require.Nil(t, result, "Should not return a list")
}

func TestLast(t *testing.T) {
	// Arrange
	input := []int{10, 14, 28}
	target := 28

	// Act
	result := generics.Last(input)

	// Assert
	require.Equal(t, result, target, "Should have the correct value")
}

func TestLastEmpty(t *testing.T) {
	// Arrange
	input := []int{}

	// Act
	result := generics.Last(input)

	// Assert
	require.Equal(t, 0, result, "Should get the default value of the type")
}

func TestLastWithContext(t *testing.T) {
	// Arrange
	testCtx := context.Background()
	input := []int{10, 30, 14, 10}
	filter := func(ctx context.Context, index int, v int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		if v >= 12 {
			return true, nil
		}
		return false, nil
	}
	target := 14

	// Act
	result, err := generics.LastWithContext(testCtx, input, filter)

	// Assert
	require.NoError(t, err, "Should not error")
	require.Equal(t, result, target, "Should have the correct value")
}

func TestLastWithContextNotFound(t *testing.T) {
	// Arrange
	testCtx := context.Background()
	input := []int{10, 14, 28}
	filter := func(ctx context.Context, index int, v int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		return false, nil
	}
	target := 0

	// Act
	result, err := generics.LastWithContext(testCtx, input, filter)

	// Assert
	require.NoError(t, err, "Should not error")
	require.Equal(t, result, target, "Should have the correct value")
}

func TestLastWithContextError(t *testing.T) {
	// Arrange
	testCtx := context.Background()
	input := []int{10, 14, 28}
	filter := func(ctx context.Context, index int, v int) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		if v >= 12 {
			return false, errors.New("boom")
		}
		return false, nil
	}

	// Act
	result, err := generics.LastWithContext(testCtx, input, filter)

	// Assert
	require.Error(t, err, "Should not error")
	require.Equal(t, result, 0, "Should have the default type value")
}
