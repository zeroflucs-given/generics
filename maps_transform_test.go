package generics_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeroflucs-given/generics"
)

func TestKeyValuesToMap(t *testing.T) {
	// Arrange
	input := []generics.KeyValuePair[string, string]{
		{
			Key:   "hello",
			Value: "world",
		},
		{
			Key:   "fizz",
			Value: "buzz",
		},
	}

	// Act
	result := generics.KeyValuesToMap(input)

	// Assert
	require.NotNil(t, result, "Should have a result")
	require.Len(t, result, 2, "Should have right number of keys")
	require.Equal(t, map[string]string{
		"hello": "world",
		"fizz":  "buzz",
	}, result, "Should have expected output")
}

func TestKeys(t *testing.T) {
	// Arrange
	input := map[string]string{
		"hello": "world",
		"fizz":  "buzz",
	}

	// Act
	result := generics.Keys(input)

	// Assert
	require.NotNil(t, result, "Should have a result")
	require.ElementsMatch(t, result, []string{"hello", "fizz"})
}

func TestMapValues(t *testing.T) {
	// Arrange
	input := map[int]string{
		1: "one",
		2: "two",
	}

	// Act
	remapped := generics.MapValues(input, func(k int, v string) bool {
		return v == "one"
	})

	// Assert
	require.Equal(t, true, remapped[1])
	require.Equal(t, false, remapped[2])
}

func TestMapValuesWithContext(t *testing.T) {
	// Arrange
	ctx := context.TODO()
	input := map[int]string{
		1: "one",
		2: "two",
	}

	// Act
	remapped, err := generics.MapValuesWithContext(ctx, input, func(ctx context.Context, k int, v string) (bool, error) {
		return v == "one", nil
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, true, remapped[1])
	require.Equal(t, false, remapped[2])
}

func TestMapValuesWithContextError(t *testing.T) {
	// Arrange
	ctx := context.TODO()
	input := map[int]string{
		1: "one",
		2: "two",
	}
	fail := errors.New("fail")

	// Act
	remapped, err := generics.MapValuesWithContext(ctx, input, func(ctx context.Context, k int, v string) (bool, error) {
		if v == "two" {
			return false, fail
		}
		return v == "one", nil
	})

	// Assert
	require.ErrorIs(t, err, fail)
	require.Nil(t, remapped)
}

func TestValues(t *testing.T) {
	// Arrange
	input := map[string]string{
		"hello": "world",
		"fizz":  "buzz",
	}

	// Act
	result := generics.Values(input)

	// Assert
	require.NotNil(t, result, "Should have a result")
	require.ElementsMatch(t, result, []string{"world", "buzz"})
}

func TestToKeyValues(t *testing.T) {
	// Arrange
	input := map[string]string{
		"hello": "world",
		"fizz":  "buzz",
	}

	// Act
	result := generics.ToKeyValues(input)

	// Assert
	require.NotNil(t, result, "Should have a result")
	require.ElementsMatch(t, result, []generics.KeyValuePair[string, string]{
		{
			Key:   "hello",
			Value: "world",
		},
		{
			Key:   "fizz",
			Value: "buzz",
		},
	})
}
