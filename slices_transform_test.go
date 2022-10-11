package generics_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeroflucs-given/generics"
)

// TestGroup checks that we can group items
func TestGroup(t *testing.T) {
	// Arrange
	mapper := func(index int, v int) string {
		if v%2 == 0 {
			return "Evens"
		} else {
			return "Odds"
		}
	}
	input := []int{1, 2, 3, 4, 5}

	// Act
	result := generics.Group(input, mapper)

	// Assert
	require.Len(t, result, 2, "Should have 2 groups")
	require.Equal(t, []int{2, 4}, result["Evens"], "Should have the right Evens")
	require.Equal(t, []int{1, 3, 5}, result["Odds"], "Should have the right Odds")
}

// TestGroupWithContext checks that we can group items
func TestGroupWithContext(t *testing.T) {
	// Arrange
	testCtx := context.Background()
	mapper := func(ctx context.Context, index int, v int) (string, error) {
		require.Equal(t, testCtx, ctx, "Should have correct context")
		if v%2 == 0 {
			return "Evens", nil
		}
		return "Odds", nil
	}
	input := []int{1, 2, 3, 4, 5}

	// Act
	result, err := generics.GroupWithContext(testCtx, input, mapper)

	// Assert
	require.NoError(t, err, "Should not error")
	require.Len(t, result, 2, "Should have 2 groups")
	require.Equal(t, []int{2, 4}, result["Evens"], "Should have the right Evens")
	require.Equal(t, []int{1, 3, 5}, result["Odds"], "Should have the right Odds")
}

// TestGroupWithContextFail checks that we can group items but stop on failure
func TestGroupWithContextFail(t *testing.T) {
	// Arrange
	testCtx := context.Background()
	mapper := func(ctx context.Context, index int, v int) (string, error) {
		require.Equal(t, testCtx, ctx, "Should have correct context")
		if v == 4 {
			return "", errors.New("boom")
		}

		if v%2 == 0 {
			return "Evens", nil
		}
		return "Odds", nil
	}
	input := []int{1, 2, 3, 4, 5}

	// Act
	result, err := generics.GroupWithContext(testCtx, input, mapper)

	// Assert
	require.Error(t, err, "Should error")
	require.Nil(t, result, "No result should be given")
}

// TestMap tests that we can map values from one type to another
func TestMap(t *testing.T) {
	// Arrange
	mapper := func(index int, v int) string {
		return fmt.Sprintf("%v", v)
	}
	input := []int{1, 2, 3, 4, 5, 6}

	// Act
	result := generics.Map(input, mapper)

	// Assert
	require.Len(t, result, 6, "Should have 6 items")
	require.Equal(t, []string{"1", "2", "3", "4", "5", "6"}, result, "Should have the right outputs")
}

// TestMapContext applies a mapper to a set with a context
func TestMapContext(t *testing.T) {
	// Arrange
	testCtx := context.TODO()
	mapper := func(ctx context.Context, index int, v int) (string, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		return fmt.Sprintf("%v", v), nil
	}
	input := []int{1, 2, 3, 4, 5, 6}

	// Act
	result, err := generics.MapWithContext(testCtx, input, mapper)

	// Assert
	require.NoError(t, err, "Should not error")
	require.Len(t, result, 6, "Should have 6 items")
	require.Equal(t, []string{"1", "2", "3", "4", "5", "6"}, result, "Should have the right outputs")
}

// TestMapContextError make sure we abort when errors occur
func TestMapContextError(t *testing.T) {
	// Arrange
	testCtx := context.TODO()
	mapper := func(ctx context.Context, index int, v int) (string, error) {
		if v == 5 {
			return "", errors.New("boom")
		}
		require.Equal(t, testCtx, ctx, "Should have right context")
		return fmt.Sprintf("%v", v), nil
	}
	input := []int{1, 2, 3, 4, 5, 6}

	// Act
	result, err := generics.MapWithContext(testCtx, input, mapper)

	// Assert
	require.Error(t, err, "Should error")
	require.Nil(t, result, "Should have no result")
}

func TestToMap(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4, 5}
	keyMapper := func(index int, k int) string {
		return fmt.Sprintf("Key%v", k%4)
	}
	valueMapper := func(index int, v int) int {
		return v * 2
	}

	// Act
	result := generics.ToMap(input, keyMapper, valueMapper)

	// Assert
	assert.Len(t, result, 4, "Should have 4 items in the dictionary")
	assert.Equal(t, result["Key1"], 10, "Should have overwritten the first value")
	assert.Equal(t, result["Key2"], 4, "Should have the right value for second item")
}

func TestToMapContext(t *testing.T) {
	// Arrange
	testCtx := context.TODO()
	input := []int{1, 2, 3, 4, 5}
	keyMapper := func(ctx context.Context, index int, k int) (string, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		return fmt.Sprintf("Key%v", k%4), nil
	}
	valueMapper := func(ctx context.Context, index int, v int) (int, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		return v * 2, nil
	}

	// Act
	result, err := generics.ToMapWithContext(testCtx, input, keyMapper, valueMapper)

	// Assert
	require.NoError(t, err, "Should not error")
	assert.Len(t, result, 4, "Should have 4 items in the dictionary")
	assert.Equal(t, result["Key1"], 10, "Should have overwritten the first value")
	assert.Equal(t, result["Key2"], 4, "Should have the right value for second item")
}

func TestToMapContextErrorKey(t *testing.T) {
	// Arrange
	testCtx := context.TODO()
	input := []int{1, 2, 3, 4, 5}
	keyMapper := func(ctx context.Context, index int, k int) (string, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		if index == 3 {
			return "", fmt.Errorf("key failed")
		}
		return fmt.Sprintf("Key%v", k%4), nil
	}
	valueMapper := func(ctx context.Context, index int, v int) (int, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		return v * 2, nil
	}

	// Act
	result, err := generics.ToMapWithContext(testCtx, input, keyMapper, valueMapper)

	// Assert
	require.Error(t, err, "Should error")
	require.Nil(t, result, "Should have no result")
}

func TestToMapContextErrorValue(t *testing.T) {
	// Arrange
	testCtx := context.TODO()
	input := []int{1, 2, 3, 4, 5}
	keyMapper := func(ctx context.Context, index int, k int) (string, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		return fmt.Sprintf("Key%v", k%4), nil
	}
	valueMapper := func(ctx context.Context, index int, v int) (int, error) {
		require.Equal(t, testCtx, ctx, "Should have right context")
		if index == 3 {
			return 0, fmt.Errorf("value failed")
		}
		return v * 2, nil
	}

	// Act
	result, err := generics.ToMapWithContext(testCtx, input, keyMapper, valueMapper)

	// Assert
	require.Error(t, err, "Should error")
	require.Nil(t, result, "Should have no result")
}

// TestCompactSlice makes sure we compact a slice so we keep only the last value per
// key in the slice
func TestCompactSlice(t *testing.T) {
	type input struct {
		K string
		V string
	}

	// Arrange
	inputData := []input{
		{
			K: "foo",
			V: "first-foo",
		},
		{
			K: "bar",
			V: "first-bar",
		},
		{
			K: "foo",
			V: "second-foo",
		},
		{
			K: "foo",
			V: "third-foo",
		},
		{
			K: "fizz",
			V: "first-fizz",
		},
		{
			K: "bar",
			V: "second-bar",
		},
	}

	// Act
	result := generics.Compact(inputData, func(i int, v input) string {
		return v.K
	})

	// Assert
	require.Len(t, result, 3, "Should have 3 items")
	require.Equal(t, "third-foo", result[0].V)
	require.Equal(t, "first-fizz", result[1].V)
	require.Equal(t, "second-bar", result[2].V)
}
