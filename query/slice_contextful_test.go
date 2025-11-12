package query_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeroflucs-given/generics/query"
)

// TestQueryableSliceWithContext just checks we can chain up long sequences of mutations
// using contexts
func TestQueryableSliceWithContext(t *testing.T) {
	// Arrange
	inputs := [][]int{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10},
		{11, 12, 13, 14, 15},
		{16, 17, 18, 19, 20},
	}
	testCtx := context.Background()

	// Act
	query := query.Slice(inputs[0]).
		WithContext(testCtx).
		Concatenate(inputs[1:]...).
		TakeWhile(func(ctx context.Context, index int, v int) (bool, error) {
			return index < 18, nil
		}).
		Skip(5).
		TakeUntil(func(ctx context.Context, index int, v int) (bool, error) {
			return v >= 15, nil
		}).
		Filter(func(ctx context.Context, index int, v int) (bool, error) {
			return v%2 == 0, nil // Evens only
		}).
		Take(3).
		Reverse()

	// Assert
	all, errAll := query.All(func(ctx context.Context, index int, v int) (bool, error) {
		return v < 12, nil
	})
	require.NoError(t, errAll, "Should not error getting .All()")
	require.True(t, all, "All items should be < 12")

	any, errAny := query.Any(func(ctx context.Context, index int, v int) (bool, error) {
		return v%2 == 1, nil
	})
	require.NoError(t, errAny, "Should not error getting .Any()")
	require.False(t, any, "Should have no odd values")

	count, errCount := query.Count()
	require.NoError(t, errCount, "Should not error getting .Count()")
	require.Equal(t, 3, count, "Should have 3 items")

	result, errSlice := query.ToSlice()
	require.NoError(t, errSlice, "Should not error getting back to slice")
	require.Equal(t, []int{10, 8, 6}, result, "Should have right items")

	first, errFirst := query.First()
	require.NoError(t, errFirst, "Should not error getting .First()")
	require.Equal(t, 10, first, "First item should be 10")

	last, errLast := query.Last()
	require.NoError(t, errLast, "Should not error getting .Last()")
	require.Equal(t, 6, last, "Last item should be 6")
}

func failFilter[T any](ctx context.Context, index int, v T) (bool, error) {
	return false, errors.New("boom")
}

func panicFilter[T any](ctx context.Context, intdex int, v T) (bool, error) {
	panic("boom")
}

// TestContextAllLazy checks that we dont' call a second filter that panics if the first filter
// is an error
func TestContextAllLazy(t *testing.T) {
	// Arrange
	ctx := context.Background()
	items := []int{1, 2, 3}
	query := query.Slice(items).WithContext(ctx)

	// Act
	result, err := query.Filter(failFilter[int]).
		All(panicFilter[int])

	// Require
	require.Error(t, err, "Should error (and not have paniced)")
	require.False(t, result, "Should have false result")
}

// TestContextAnyLazy checks that we dont' call a second filter that panics if the first filter
// is an error
func TestContextAnyLazy(t *testing.T) {
	// Arrange
	ctx := context.Background()
	items := []int{1, 2, 3}
	query := query.Slice(items).WithContext(ctx)

	// Act
	result, err := query.Filter(failFilter[int]).
		Any(panicFilter[int])

	// Require
	require.Error(t, err, "Should error (and not have paniced)")
	require.False(t, result, "Should have false result")
}

// TestContextConcatenateLazy checks we don't clone data when being lazy
func TestContextConcatenateLazy(t *testing.T) {
	// Arrange
	ctx := context.Background()
	items := []int{1, 2, 3}
	query := query.Slice(items).WithContext(ctx)

	// Act
	original := query.Filter(failFilter[int])
	second := original.Concatenate(nil)
	result, err := second.ToSlice()

	// Require
	require.Error(t, err, "Should error (and not have paniced)")
	require.Nil(t, result, "Should have no result")
}

// TestContextCountLazy checks we don't call filters on a count when we've already failed
func TestContextCountLazy(t *testing.T) {
	// Arrange
	ctx := context.Background()
	items := []int{1, 2, 3}
	query := query.Slice(items).WithContext(ctx)

	// Act
	result, err := query.Filter(failFilter[int]).
		Count(panicFilter[int])

	// Require
	require.Error(t, err, "Should error (and not have paniced)")
	require.Zero(t, result, "Should have 0 result")
}

// TestContextFilterLazy checks we don't call filters on a filter when we've already failed
func TestContextFilterLazy(t *testing.T) {
	// Arrange
	ctx := context.Background()
	items := []int{1, 2, 3}
	query := query.Slice(items).WithContext(ctx)

	// Act
	result, err := query.Filter(failFilter[int]).
		Filter(panicFilter[int]).
		ToSlice()

	// Require
	require.Error(t, err, "Should error (and not have paniced)")
	require.Nil(t, result, "Should have nil result")
}

// TestContextFirstLazy checks we don't call filters on a first when we've already failed
func TestContextFirstLazy(t *testing.T) {
	// Arrange
	ctx := context.Background()
	items := []int{1, 2, 3}
	query := query.Slice(items).WithContext(ctx)

	// Act
	result, err := query.Filter(failFilter[int]).
		First(panicFilter[int])

	// Require
	require.Error(t, err, "Should error (and not have paniced)")
	require.Zero(t, result, "Should have nil result")
}

// TestContextFirstLazy checks we don't call filters on a first when we've already failed
func TestContextLastLazy(t *testing.T) {
	// Arrange
	ctx := context.Background()
	items := []int{1, 2, 3}
	query := query.Slice(items).WithContext(ctx)

	// Act
	result, err := query.Filter(failFilter[int]).
		Last(panicFilter[int])

	// Require
	require.Error(t, err, "Should error (and not have paniced)")
	require.Zero(t, result, "Should have nil result")
}

// TestContextReverseLazy checks we don't clone data when being lazy
func TestContextReverseLazy(t *testing.T) {
	// Arrange
	ctx := context.Background()
	items := []int{1, 2, 3}
	query := query.Slice(items).WithContext(ctx)

	// Act
	original := query.Filter(failFilter[int])
	second := original.Reverse()
	result, err := second.ToSlice()

	// Require
	require.Error(t, err, "Should error (and not have paniced)")
	require.Nil(t, result, "Should have no result")
	// require.Equal(t, original, &second, "Should be reference equal")
}

// TestContextSkipLazy checks we don't skip data when being lazy
func TestContextSkipLazy(t *testing.T) {
	// Arrange
	ctx := context.Background()
	items := []int{1, 2, 3}
	query := query.Slice(items).WithContext(ctx)

	// Act
	original := query.Filter(failFilter[int])
	second := original.Skip(10)
	result, err := second.ToSlice()

	// Require
	require.Error(t, err, "Should error (and not have paniced)")
	require.Nil(t, result, "Should have no result")
	// require.Equal(t, original, &second, "Should be reference equal")
}

// TestContextTakeLazy checks we don't skip data when being lazy
func TestContextTakeLazy(t *testing.T) {
	// Arrange
	ctx := context.Background()
	items := []int{1, 2, 3}
	query := query.Slice(items).WithContext(ctx)

	// Act
	original := query.Filter(failFilter[int])
	second := original.Take(10)
	result, err := second.ToSlice()

	// Require
	require.Error(t, err, "Should error (and not have paniced)")
	require.Nil(t, result, "Should have no result")
	// require.Equal(t, original, &second, "Should be reference equal")
}

// TestContextTakeUntilLazy checks we don't take data when being lazy
func TestContextTakeUntilLazy(t *testing.T) {
	// Arrange
	ctx := context.Background()
	items := []int{1, 2, 3}
	query := query.Slice(items).WithContext(ctx)

	// Act
	original := query.Filter(failFilter[int])
	second := original.TakeUntil(panicFilter[int])
	result, err := second.ToSlice()

	// Require
	require.Error(t, err, "Should error (and not have paniced)")
	require.Nil(t, result, "Should have no result")
	// require.Equal(t, original, &second, "Should be reference equal")
}

// TestContextTakeWhileLazy checks we don't take data when being lazy
func TestContextTakeWhileLazy(t *testing.T) {
	// Arrange
	ctx := context.Background()
	items := []int{1, 2, 3}
	query := query.Slice(items).WithContext(ctx)

	// Act
	original := query.Filter(failFilter[int])
	second := original.TakeWhile(panicFilter[int])
	result, err := second.ToSlice()

	// Require
	require.Error(t, err, "Should error (and not have paniced)")
	require.Nil(t, result, "Should have no result")
	// require.Equal(t, original, &second, "Should be reference equal")
}
