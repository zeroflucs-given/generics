package filtering_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeroflucs-given/generics/filtering"
)

// TestAnd checks we call all filters
func TestAnd(t *testing.T) {
	// Arrange
	var first, second bool
	filterFirst := func(index int, v string) bool {
		first = true
		return true
	}
	filterSecond := func(index int, v string) bool {
		second = true
		return true
	}

	// Act
	wrapped := filtering.And(filterFirst, filterSecond)
	result := wrapped(1, "hello")

	// Assert
	require.True(t, result, "Should have right result")
	require.True(t, first, "Should call the first filter")
	require.True(t, second, "Should call the second filter")
}

// TestAndFail checks we call filters in a lazy fashion
func TestAndFail(t *testing.T) {
	// Arrange
	var first, second bool
	filterFirst := func(index int, v string) bool {
		first = true
		return false
	}
	filterSecond := func(index int, v string) bool {
		second = true
		return true
	}

	// Act
	wrapped := filtering.And(filterFirst, filterSecond)
	result := wrapped(1, "hello")

	// Assert
	require.False(t, result, "Should have right result")
	require.True(t, first, "Should call the first filter")
	require.False(t, second, "Should not call the second filter")
}

// TestAndWithContext checks we call all filters
func TestAndWithContext(t *testing.T) {
	// Arrange
	var first, second bool
	testCtx := context.Background()
	filterFirst := func(ctx context.Context, index int, v string) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context if called")
		first = true
		return true, nil
	}
	filterSecond := func(ctx context.Context, index int, v string) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context if called")
		second = true
		return true, nil
	}

	// Act
	wrapped := filtering.AndWithContext(filterFirst, filterSecond)
	result, err := wrapped(testCtx, 1, "hello")

	// Assert
	require.NoError(t, err, "Should not error")
	require.True(t, result, "Should have right result")
	require.True(t, first, "Should call the first filter")
	require.True(t, second, "Should call the second filter")
}

// TestAndWithContextFail checks we call filters in a lazy fashion
func TestAndWithContextFail(t *testing.T) {
	// Arrange
	var first, second bool
	testCtx := context.Background()
	filterFirst := func(ctx context.Context, index int, v string) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context if called")
		first = true
		return false, nil
	}
	filterSecond := func(ctx context.Context, index int, v string) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context if called")
		second = true
		return true, nil
	}

	// Act
	wrapped := filtering.AndWithContext(filterFirst, filterSecond)
	result, err := wrapped(testCtx, 1, "hello")

	// Assert
	require.NoError(t, err, "Should not error")
	require.False(t, result, "Should have right result")
	require.True(t, first, "Should call the first filter")
	require.False(t, second, "Should not call the second filter")
}

// TestAndWithContextError checks we call filters in a lazy fashion and propegate errors
func TestAndWithContextError(t *testing.T) {
	// Arrange
	var first, second bool
	testCtx := context.Background()
	filterFirst := func(ctx context.Context, index int, v string) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context if called")
		first = true
		return false, errors.New("boom")
	}
	filterSecond := func(ctx context.Context, index int, v string) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context if called")
		second = true
		return true, nil
	}

	// Act
	wrapped := filtering.AndWithContext(filterFirst, filterSecond)
	result, err := wrapped(testCtx, 1, "hello")

	// Assert
	require.Error(t, err, "Should error")
	require.False(t, result, "Should have right result")
	require.True(t, first, "Should call the first filter")
	require.False(t, second, "Should not call the second filter")
}

// TestOr checks we call all filters until we need
func TestOr(t *testing.T) {
	// Arrange
	var first, second bool
	filterFirst := func(index int, v string) bool {
		first = true
		return false // Fall through to second
	}
	filterSecond := func(index int, v string) bool {
		second = true
		return true
	}

	// Act
	wrapped := filtering.Or(filterFirst, filterSecond)
	result := wrapped(1, "hello")

	// Assert
	require.True(t, result, "Should have right result")
	require.True(t, first, "Should call the first filter")
	require.True(t, second, "Should call the second filter")
}

// TestOrLazy checks we call all filters in a lazy way
func TestOrLazy(t *testing.T) {
	// Arrange
	var first, second bool
	filterFirst := func(index int, v string) bool {
		first = true
		return true // Don't fall through to second
	}
	filterSecond := func(index int, v string) bool {
		second = true
		return true
	}

	// Act
	wrapped := filtering.Or(filterFirst, filterSecond)
	result := wrapped(1, "hello")

	// Assert
	require.True(t, result, "Should have right result")
	require.True(t, first, "Should call the first filter")
	require.False(t, second, "Should call the second filter")
}

// TestOrFail checks we call filters until we get a result
func TestOrFail(t *testing.T) {
	// Arrange
	var first, second bool
	filterFirst := func(index int, v string) bool {
		first = true
		return false
	}
	filterSecond := func(index int, v string) bool {
		second = true
		return false
	}

	// Act
	wrapped := filtering.Or(filterFirst, filterSecond)
	result := wrapped(1, "hello")

	// Assert
	require.False(t, result, "Should have right result")
	require.True(t, first, "Should call the first filter")
	require.True(t, second, "Should not call the second filter")
}

// TestOrWithContext checks we call all filters until we get one
func TestOrWithContext(t *testing.T) {
	// Arrange
	var first, second bool
	testCtx := context.Background()
	filterFirst := func(ctx context.Context, index int, v string) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context if called")
		first = true
		return false, nil // Will cause second to run
	}
	filterSecond := func(ctx context.Context, index int, v string) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context if called")
		second = true
		return true, nil
	}

	// Act
	wrapped := filtering.OrWithContext(filterFirst, filterSecond)
	result, err := wrapped(testCtx, 1, "hello")

	// Assert
	require.NoError(t, err, "Should not error")
	require.True(t, result, "Should have right result")
	require.True(t, first, "Should call the first filter")
	require.True(t, second, "Should call the second filter")
}

// TestOrWithContextFail checks we call filters until we need to
func TestOrWithContextFail(t *testing.T) {
	// Arrange
	var first, second bool
	testCtx := context.Background()
	filterFirst := func(ctx context.Context, index int, v string) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context if called")
		first = true
		return false, nil
	}
	filterSecond := func(ctx context.Context, index int, v string) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context if called")
		second = true
		return false, nil
	}

	// Act
	wrapped := filtering.OrWithContext(filterFirst, filterSecond)
	result, err := wrapped(testCtx, 1, "hello")

	// Assert
	require.NoError(t, err, "Should not error")
	require.False(t, result, "Should have right result")
	require.True(t, first, "Should call the first filter")
	require.True(t, second, "Should call the second filter")
}

// TestOrWithContextError checks we call filters in a lazy fashion and propegate errors
func TestOrWithContextError(t *testing.T) {
	// Arrange
	var first, second bool
	testCtx := context.Background()
	filterFirst := func(ctx context.Context, index int, v string) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context if called")
		first = true
		return false, errors.New("boom")
	}
	filterSecond := func(ctx context.Context, index int, v string) (bool, error) {
		require.Equal(t, testCtx, ctx, "Should have right context if called")
		second = true
		return true, nil
	}

	// Act
	wrapped := filtering.OrWithContext(filterFirst, filterSecond)
	result, err := wrapped(testCtx, 1, "hello")

	// Assert
	require.Error(t, err, "Should error")
	require.False(t, result, "Should have right result")
	require.True(t, first, "Should call the first filter")
	require.False(t, second, "Should not call the second filter")
}

func TestNot(t *testing.T) {
	// Arrange
	base := filtering.True[string]
	wrapped := filtering.Not(base)

	// Act
	result := wrapped(1, "hello")

	// Assert
	require.False(t, result, "Should have correct result")
}

func TestNotWithContext(t *testing.T) {
	// Arrange
	base := filtering.TrueWithContext[string]
	wrapped := filtering.NotWithContext(base)
	testCtx := context.Background()

	// Act
	result, err := wrapped(testCtx, 1, "hello")

	// Assert
	require.NoError(t, err, "Should not error")
	require.False(t, result, "Should have correct result")
}

func TestNotWithContextError(t *testing.T) {
	// Arrange
	base := func(ctx context.Context, index int, v string) (bool, error) {
		return false, errors.New("boom")
	}
	wrapped := filtering.NotWithContext(base)
	testCtx := context.Background()

	// Act
	result, err := wrapped(testCtx, 1, "hello")

	// Assert
	require.Error(t, err, "Should propegate error")
	require.False(t, result, "Should have correct result")
}
